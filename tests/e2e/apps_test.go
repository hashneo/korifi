package e2e_test

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Apps", func() {
	var (
		space1GUID string
		space1Name string
		appGUID    string
		resp       *resty.Response
	)

	BeforeEach(func() {
		space1Name = generateGUID("space1")
		space1GUID = createSpace(space1Name, commonTestOrgGUID)
	})

	AfterEach(func() {
		deleteSpace(space1GUID)
	})

	Describe("List apps", func() {
		var (
			space2GUID, space3GUID       string
			app1GUID, app2GUID, app3GUID string
			app4GUID, app5GUID, app6GUID string
			result                       resourceList[resource]
		)

		BeforeEach(func() {
			space2GUID = createSpace(generateGUID("space2"), commonTestOrgGUID)
			space3GUID = createSpace(generateGUID("space3"), commonTestOrgGUID)

			createSpaceRole("space_developer", certUserName, space1GUID)
			createSpaceRole("space_developer", certUserName, space3GUID)

			app1GUID = createApp(space1GUID, generateGUID("app1"))
			app2GUID = createApp(space1GUID, generateGUID("app2"))
			app3GUID = createApp(space2GUID, generateGUID("app3"))
			app4GUID = createApp(space2GUID, generateGUID("app4"))
			app5GUID = createApp(space3GUID, generateGUID("app5"))
			app6GUID = createApp(space3GUID, generateGUID("app6"))
		})

		AfterEach(func() {
			deleteSpace(space2GUID)
			deleteSpace(space3GUID)
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetResult(&result).Get("/v3/apps")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns apps only in authorized spaces", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusOK))

			Expect(result.Resources).To(ContainElements(
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app1GUID)}),
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app2GUID)}),
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app5GUID)}),
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app6GUID)}),
			))

			Expect(result.Resources).ToNot(ContainElements(
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app3GUID)}),
				MatchFields(IgnoreExtras, Fields{"GUID": Equal(app4GUID)}),
			))
		})
	})

	Describe("Create an app", func() {
		var appName string

		BeforeEach(func() {
			appName = generateGUID("app")
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetBody(appResource{
				resource: resource{
					Name: appName,
					Relationships: relationships{
						"space": {
							Data: resource{
								GUID: space1GUID,
							},
						},
					},
				},
			}).Post("/v3/apps")
			Expect(err).NotTo(HaveOccurred())
		})

		When("the user has space developer role in the space", func() {
			BeforeEach(func() {
				createSpaceRole("space_developer", certUserName, space1GUID)
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusCreated))
			})

			When("the app name already exists in the space", func() {
				BeforeEach(func() {
					createApp(space1GUID, appName)
				})

				It("returns an unprocessable entity error", func() {
					Expect(resp).To(HaveRestyStatusCode(http.StatusUnprocessableEntity))
					Expect(resp).To(HaveRestyBody(ContainSubstring("CF-UniquenessError")))
				})
			})
		})

		When("the user cannot create apps in the space", func() {
			BeforeEach(func() {
				createSpaceRole("space_manager", certUserName, space1GUID)
			})

			It("fails", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusForbidden))
				Expect(resp).To(HaveRestyBody(ContainSubstring("CF-NotAuthorized")))
			})
		})
	})

	Describe("Delete an app", func() {
		var (
			err         error
			pkgGUID     string
			buildGUID   string
			processType string
		)

		BeforeEach(func() {
			createSpaceRole("space_developer", certUserName, space1GUID)
			appName := generateGUID("app")
			appGUID = createApp(space1GUID, appName)
			pkgGUID = createPackage(appGUID)
			buildGUID = createBuild(pkgGUID)
			processType = "web"

			// Applying a manifest will create any processes defined within the manifest
			manifest := manifestResource{
				Version: 1,
				Applications: []applicationResource{{
					Name: appName,
					Processes: []manifestApplicationProcessResource{{
						Type: processType,
					}},
				}},
			}

			applySpaceManifest(manifest, space1GUID)
			var processResult processResource
			resp, err = certClient.R().
				SetResult(&processResult).
				Get("/v3/apps/" + appGUID + "/processes/" + processType)
			Expect(err).NotTo(HaveOccurred())
			Expect(processResult.Type).To(Equal(processType))

			resp, err = certClient.R().Delete("/v3/apps/" + appGUID)
		})

		It("succeeds with a job redirect", func() {
			Expect(resp).To(SatisfyAll(
				HaveRestyStatusCode(http.StatusAccepted),
				HaveRestyHeaderWithValue("Location", HaveSuffix("/v3/jobs/app.delete~"+appGUID)),
			))

			jobURL := resp.Header().Get("Location")
			Eventually(func(g Gomega) {
				resp, err = certClient.R().Get(jobURL)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(string(resp.Body())).To(ContainSubstring("COMPLETE"))
			}).Should(Succeed())
		})

		It("successfully deletes the app and associated child resources", func() {
			var result resource
			Eventually(func(g Gomega) {
				resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID)
				g.Expect(resp).To(HaveRestyStatusCode(http.StatusNotFound))
				g.Expect(err).NotTo(HaveOccurred())
			}).Should(Succeed())

			Eventually(func(g Gomega) {
				resp, err = certClient.R().SetResult(&result).Get("/v3/packages/" + pkgGUID)
				g.Expect(resp).To(HaveRestyStatusCode(http.StatusNotFound))
				g.Expect(err).NotTo(HaveOccurred())
			}).Should(Succeed())

			Eventually(func(g Gomega) {
				resp, err = certClient.R().SetResult(&result).Get("/v3/builds/" + buildGUID)
				g.Expect(resp).To(HaveRestyStatusCode(http.StatusNotFound))
				g.Expect(err).NotTo(HaveOccurred())
			}).Should(Succeed())

			Eventually(func(g Gomega) {
				resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/processes/" + processType)
				g.Expect(resp).To(HaveRestyStatusCode(http.StatusNotFound))
				g.Expect(err).NotTo(HaveOccurred())
			}).Should(Succeed())
		})
	})

	Describe("Fetch an app", func() {
		var result resource

		BeforeEach(func() {
			createSpaceRole("space_developer", certUserName, space1GUID)
			appGUID = createApp(space1GUID, generateGUID("app1"))
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can fetch the app", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
			Expect(result.GUID).To(Equal(appGUID))
		})
	})

	Describe("List app processes", func() {
		var result resourceList[typedResource]

		BeforeEach(func() {
			createSpaceRole("space_developer", certUserName, space1GUID)
			appGUID = createApp(space1GUID, generateGUID("app"))
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/processes")
			Expect(err).NotTo(HaveOccurred())
		})

		It("successfully lists the default web process", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
			Expect(result.Resources).To(HaveLen(1))
			Expect(result.Resources[0].Type).To(Equal("web"))
		})
	})

	Describe("Get app process by type", func() {
		var (
			result      resource
			processGUID string
		)

		BeforeEach(func() {
			appGUID, _ = pushTestApp(space1GUID, appBitsFile)
			processGUID = getProcess(appGUID, "web").GUID
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/processes/web")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a not-found error to users with no space access", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusNotFound))
			Expect(resp).To(HaveRestyBody(ContainSubstring("App not found")))
		})

		When("the user is a space developer", func() {
			BeforeEach(func() {
				createSpaceRole("space_developer", certUserName, space1GUID)
			})

			It("successfully returns the process", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.GUID).To(Equal(processGUID))
			})
		})
	})

	Describe("List app routes", func() {
		var result resourceList[resource]

		BeforeEach(func() {
			createSpaceRole("space_developer", certUserName, space1GUID)
			appGUID = createApp(space1GUID, generateGUID("app"))
		})

		JustBeforeEach(func() {
			var err error
			resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/routes")
			Expect(err).NotTo(HaveOccurred())
		})

		It("successfully lists the empty set of routes", func() {
			Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
			Expect(result.Resources).To(BeEmpty())
		})
	})

	Describe("Built apps", func() {
		var (
			pkgGUID   string
			buildGUID string
			result    resource
		)

		BeforeEach(func() {
			appGUID = createApp(space1GUID, generateGUID("app"))
			pkgGUID = createPackage(appGUID)
			uploadTestApp(pkgGUID, appBitsFile)
			buildGUID = createBuild(pkgGUID)
			waitForDroplet(buildGUID)

			createSpaceRole("space_developer", certUserName, space1GUID)
		})

		Describe("Get app current droplet", func() {
			BeforeEach(func() {
				setCurrentDroplet(appGUID, buildGUID)
			})

			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/droplets/current")
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.GUID).To(Equal(buildGUID))
			})
		})

		Describe("Set app current droplet", func() {
			type currentDropletResource struct {
				Data resource `json:"data"`
			}

			var result currentDropletResource

			JustBeforeEach(func() {
				var err error
				var body currentDropletResource
				body.Data.GUID = buildGUID

				resp, err = certClient.R().
					SetBody(body).
					SetResult(&result).
					Patch("/v3/apps/" + appGUID + "/relationships/current_droplet")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns 200", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.Data.GUID).To(Equal(buildGUID))
			})
		})

		Describe("Start an app", func() {
			var result appResource

			BeforeEach(func() {
				setCurrentDroplet(appGUID, buildGUID)
			})

			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().SetResult(&result).Post("/v3/apps/" + appGUID + "/actions/start")
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.State).To(Equal("STARTED"))
			})
		})

		Describe("Restart an app", func() {
			var result map[string]interface{}

			BeforeEach(func() {
				setCurrentDroplet(appGUID, buildGUID)
				var err error
				resp, err = certClient.R().SetResult(&result).Post("/v3/apps/" + appGUID + "/actions/start")
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result).To(HaveKeyWithValue("state", "STARTED"))
			})

			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().SetResult(&result).Post("/v3/apps/" + appGUID + "/actions/restart")
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result).To(HaveKeyWithValue("state", "STARTED"))
			})

			It("sets the app rev to 1", func() {
				Eventually(func(g Gomega) {
					var err error
					resp, err = certClient.R().
						SetResult(&result).
						Get("/v3/apps/" + appGUID)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
					g.Expect(result).To(HaveKeyWithValue("metadata", HaveKeyWithValue("annotations", HaveKeyWithValue("korifi.cloudfoundry.org/app-rev", "1"))))
				}).Should(Succeed())
			})

			When("the app is restarted again", func() {
				JustBeforeEach(func() {
					var err error
					resp, err = certClient.R().SetResult(&result).Post("/v3/apps/" + appGUID + "/actions/restart")
					Expect(err).NotTo(HaveOccurred())
				})

				It("sets the app rev to 2", func() {
					Eventually(func(g Gomega) {
						var err error
						resp, err = certClient.R().
							SetResult(&result).
							Get("/v3/apps/" + appGUID)
						g.Expect(err).NotTo(HaveOccurred())
						g.Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
						g.Expect(result).To(HaveKeyWithValue("metadata", HaveKeyWithValue("annotations", HaveKeyWithValue("korifi.cloudfoundry.org/app-rev", "2"))))
					}).Should(Succeed())
				})
			})

			When("app environment has been changed", func() {
				BeforeEach(func() {
					setEnv(appGUID, map[string]interface{}{
						"foo": "var",
					})
				})

				It("sets the new env var to the app environment", func() {
					Expect(getEnv(appGUID)).To(HaveKeyWithValue("environment_variables", HaveKeyWithValue("foo", "var")))
				})
			})
		})

		Describe("Stop an app", func() {
			var result appResource

			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().SetResult(&result).Post("/v3/apps/" + appGUID + "/actions/stop")
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result.State).To(Equal("STOPPED"))
			})
		})
	})

	Describe("Running Apps", func() {
		var (
			appName     string
			processGUID string
		)

		BeforeEach(func() {
			appGUID, appName = pushTestApp(space1GUID, doraBitsFile)
			processGUID = getProcess(appGUID, "web").GUID
		})

		It("Runs with gettable /env.json endpoint", func() {
			body := curlApp(appGUID, "/env.json")

			env := map[string]string{}
			err := json.Unmarshal(body, &env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env).To(HaveKeyWithValue("VCAP_SERVICES", Not(BeEmpty())))
			Expect(env).To(HaveKeyWithValue("VCAP_APPLICATION", Not(BeEmpty())))
		})

		When("the app is re-pushed with different code", func() {
			BeforeEach(func() {
				body := curlApp(appGUID, "")
				Expect(body).To(ContainSubstring("Hi, I'm Dora!"))
				Expect(pushTestAppWithName(space1GUID, loggingAppBitsFile, appName)).To(Equal(appGUID))
			})

			It("it returns a different endpoint result", func() {
				body := curlApp(appGUID, "")
				Expect(body).To(ContainSubstring("hello-world from a node app!"))
			})
		})

		Describe("Scale a process", func() {
			var result responseResource
			var errResp cfErrs
			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().
					SetBody(scaleResource{Instances: 2}).
					SetError(&errResp).
					SetResult(&result).
					Post("/v3/apps/" + appGUID + "/processes/web/actions/scale")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns not found for users with no role in the space", func() {
				expectNotFoundError(resp, errResp, "App")
			})

			When("the user is a space manager", func() {
				BeforeEach(func() {
					createSpaceRole("space_manager", certUserName, space1GUID)
				})

				It("returns forbidden", func() {
					Expect(resp).To(HaveRestyStatusCode(http.StatusForbidden))
				})
			})

			When("the user is a space developer", func() {
				BeforeEach(func() {
					createSpaceRole("space_developer", certUserName, space1GUID)
				})

				It("succeeds, and returns the process", func() {
					Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
					Expect(result.GUID).To(Equal(processGUID))
				})
			})
		})

		Describe("Stop the app", func() {
			var result map[string]interface{}
			var errResp cfErrs

			BeforeEach(func() {
				createSpaceRole("space_developer", certUserName, space1GUID)
			})

			JustBeforeEach(func() {
				var err error
				resp, err = certClient.R().
					SetError(&errResp).
					SetResult(&result).
					Post("/v3/apps/" + appGUID + "/actions/stop")
				Expect(err).NotTo(HaveOccurred())
			})

			It("succeeds", func() {
				Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
				Expect(result).To(HaveKeyWithValue("state", "STOPPED"))
			})

			It("eventually increments the app-rev annotation", func() {
				Eventually(func(g Gomega) {
					resp, err := certClient.R().
						SetResult(&result).
						Get("/v3/apps/" + appGUID)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
					Expect(result).To(HaveKeyWithValue("metadata", HaveKeyWithValue("annotations", HaveKeyWithValue("korifi.cloudfoundry.org/app-rev", SatisfyAny(
						Equal("0"),
						Equal("1"),
					)))))
					g.Expect(result).To(HaveKeyWithValue("metadata", HaveKeyWithValue("annotations", HaveKeyWithValue("korifi.cloudfoundry.org/app-rev", "1"))))
				}).Should(Succeed())

				Consistently(func(g Gomega) {
					resp, err := certClient.R().
						SetResult(&result).
						Get("/v3/apps/" + appGUID)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
					g.Expect(result).To(HaveKeyWithValue("metadata", HaveKeyWithValue("annotations", HaveKeyWithValue("korifi.cloudfoundry.org/app-rev", "1"))))
				}, "2s", "200ms").Should(Succeed())
			})
		})
	})

	Describe("Fetch app env", func() {
		var (
			result                      map[string]interface{}
			instanceGUID, instanceGUID2 string
			instanceName, instanceName2 string
			bindingGUID, bindingGUID2   string
			appName                     string
		)

		BeforeEach(func() {
			createSpaceRole("space_developer", certUserName, space1GUID)
			appName = generateGUID("app1")
			appGUID = createApp(space1GUID, appName)
			setEnv(appGUID, map[string]interface{}{
				"foo": "var",
			})
			instanceName = generateGUID("service-instance")
			instanceGUID = createServiceInstance(space1GUID, instanceName)
			bindingGUID = createServiceBinding(appGUID, instanceGUID)
			instanceName2 = generateGUID("service-instance")
			instanceGUID2 = createServiceInstance(space1GUID, instanceName2)
			bindingGUID2 = createServiceBinding(appGUID, instanceGUID2)
		})

		JustBeforeEach(func() {
			Eventually(func(g Gomega) {
				resp, err := certClient.R().
					SetResult(&result).
					Get("/v3/apps/" + appGUID + "/env")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(resp).To(HaveRestyStatusCode(http.StatusOK))
			}).Should(Succeed())
		})

		It("returns the app environment", func() {
			Eventually(func(g Gomega) {
				_, err := certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/env")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(result).To(HaveKeyWithValue("environment_variables", HaveKeyWithValue("foo", "var")))
				g.Expect(result).To(SatisfyAll(
					HaveKeyWithValue("system_env_json",
						HaveKeyWithValue("VCAP_SERVICES",
							HaveKeyWithValue("user-provided", ConsistOf(
								map[string]interface{}{
									"syslog_drain_url": nil,
									"tags":             []interface{}{},
									"instance_name":    instanceName,
									"binding_guid":     bindingGUID,
									"credentials": map[string]interface{}{
										"type": "user-provided",
									},
									"volume_mounts": []interface{}{},
									"label":         "user-provided",
									"name":          instanceName,
									"instance_guid": instanceGUID,
									"binding_name":  nil,
								},
								map[string]interface{}{
									"syslog_drain_url": nil,
									"tags":             []interface{}{},
									"instance_name":    instanceName2,
									"binding_guid":     bindingGUID2,
									"credentials": map[string]interface{}{
										"type": "user-provided",
									},
									"volume_mounts": []interface{}{},
									"label":         "user-provided",
									"name":          instanceName2,
									"instance_guid": instanceGUID2,
									"binding_name":  nil,
								}),
							),
						),
					),
					HaveKeyWithValue("application_env_json",
						HaveKeyWithValue("VCAP_APPLICATION", SatisfyAll(
							HaveKeyWithValue("application_id", appGUID),
							HaveKeyWithValue("application_name", appName),
							// HaveKeyWithValue("cf_api", apiServerRoot),
							HaveKeyWithValue("name", appName),
							HaveKeyWithValue("organization_id", commonTestOrgGUID),
							HaveKeyWithValue("organization_name", commonTestOrgName),
							HaveKeyWithValue("space_id", space1GUID),
							HaveKeyWithValue("space_name", space1Name),
						)),
					),
				))
			}).Should(Succeed())
		})

		When("a deleted service binding changes the app env", func() {
			BeforeEach(func() {
				_, err := certClient.R().Delete("/v3/service_credential_bindings/" + bindingGUID2)
				Expect(err).To(Succeed())
			})

			It("returns the updated app environment", func() {
				Eventually(func(g Gomega) {
					_, err := certClient.R().SetResult(&result).Get("/v3/apps/" + appGUID + "/env")
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(result).To(HaveKeyWithValue("environment_variables", HaveKeyWithValue("foo", "var")))
					g.Expect(result).To(HaveKeyWithValue("system_env_json", HaveKeyWithValue("VCAP_SERVICES", map[string]interface{}{
						"user-provided": []interface{}{
							map[string]interface{}{
								"syslog_drain_url": nil,
								"tags":             []interface{}{},
								"instance_name":    instanceName,
								"binding_guid":     bindingGUID,
								"credentials": map[string]interface{}{
									"type": "user-provided",
								},
								"volume_mounts": []interface{}{},
								"label":         "user-provided",
								"name":          instanceName,
								"instance_guid": instanceGUID,
								"binding_name":  nil,
							},
						},
					})))
				}).Should(Succeed())
			})
		})
	})
})
