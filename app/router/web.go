package router

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudness-io/cloudness/app/auth/authn"
	"github.com/cloudness-io/cloudness/app/controller/application"
	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/deployment"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/app/controller/favorite"
	"github.com/cloudness-io/cloudness/app/controller/githubapp"
	"github.com/cloudness-io/cloudness/app/controller/gitpublic"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/logs"
	"github.com/cloudness-io/cloudness/app/controller/project"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/controller/volume"
	"github.com/cloudness-io/cloudness/app/middleware/address"
	"github.com/cloudness-io/cloudness/app/middleware/audit"
	middlewareauthn "github.com/cloudness-io/cloudness/app/middleware/authn"
	"github.com/cloudness-io/cloudness/app/middleware/hx"
	middlewareinject "github.com/cloudness-io/cloudness/app/middleware/inject"
	"github.com/cloudness-io/cloudness/app/middleware/logging"
	middlewarenav "github.com/cloudness-io/cloudness/app/middleware/nav"
	"github.com/cloudness-io/cloudness/app/middleware/nocache"
	middlewarerestrict "github.com/cloudness-io/cloudness/app/middleware/restrict"
	"github.com/cloudness-io/cloudness/app/middleware/url"
	"github.com/cloudness-io/cloudness/app/request"
	accounthandler "github.com/cloudness-io/cloudness/app/web/handler/account"
	handlerapplication "github.com/cloudness-io/cloudness/app/web/handler/application"
	authhandler "github.com/cloudness-io/cloudness/app/web/handler/auth"
	handlercreate "github.com/cloudness-io/cloudness/app/web/handler/create"
	handlerdeployment "github.com/cloudness-io/cloudness/app/web/handler/deployment"
	handlerenvironment "github.com/cloudness-io/cloudness/app/web/handler/environment"
	handlerfavorite "github.com/cloudness-io/cloudness/app/web/handler/favorite"
	handlerinstance "github.com/cloudness-io/cloudness/app/web/handler/instance"
	handlerLogs "github.com/cloudness-io/cloudness/app/web/handler/logs"
	handlerproject "github.com/cloudness-io/cloudness/app/web/handler/project"
	handlerserver "github.com/cloudness-io/cloudness/app/web/handler/server"
	handlersource "github.com/cloudness-io/cloudness/app/web/handler/source"
	handlertenant "github.com/cloudness-io/cloudness/app/web/handler/tenant"
	handlervariable "github.com/cloudness-io/cloudness/app/web/handler/variable"
	handlervolume "github.com/cloudness-io/cloudness/app/web/handler/volume"
	webhookhandler "github.com/cloudness-io/cloudness/app/web/handler/webhook"
	"github.com/cloudness-io/cloudness/app/web/public"
	"github.com/cloudness-io/cloudness/app/web/render"
	"github.com/cloudness-io/cloudness/app/web/views/dto"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
)

// WebHandler is an abstraction of an http handler that handles web calls.
type WebHandler interface {
	http.Handler
}

// NewWebHandler returns a new WebHandler.
func NewWebHandler(
	appCtx context.Context,
	config *types.Config,
	authenticator authn.Authenticator,
	instanceCtrl *instance.Controller,
	serverCtrl *server.Controller,
	userCtrl *user.Controller,
	tenantCtrl *tenant.Controller,
	projectCtrl *project.Controller,
	envCtrl *environment.Controller,
	authCtrl *auth.Controller,
	ghAppCtrl *githubapp.Controller,
	gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller,
	varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller,
	volumeCtrl *volume.Controller,
	templCtrl *template.Controller,
	favCtrl *favorite.Controller,
) WebHandler {
	// Use go-chi router for inner routing.
	r := chi.NewRouter()

	// Apply common web middleware.
	r.Use(nocache.NoCache)
	r.Use(middleware.Recoverer)

	// configure logging middleware.
	r.Use(hlog.URLHandler("http.url"))
	r.Use(hlog.MethodHandler("http.method"))
	r.Use(logging.HLogRequestIDHandler())
	r.Use(logging.HLogAccessLogHandler())
	r.Use(address.Handler("", ""))

	// configure cors middleware
	// r.Use(corsHandler(config))

	r.Use(audit.Middleware())

	// configure hx middleware
	r.Use(hx.PopulateHxIndidcator())
	r.Use(hx.PopulateHxCallerUrl())

	r.Use(url.PopulateCurrentUrl())

	r.Use(middlewareinject.PopulateTargetElemet())

	// serve static files
	if config.Environment == "local" {
		r.Handle("/public/*", disableCache(staticDev()))
	} else {
		r.Handle("/public/*", enableCache(staticProd()))
	}

	r.Route("/", func(r chi.Router) {
		r.Use(middlewareinject.InjectInstance(instanceCtrl))
		//special methods that don't require authentication
		setupAccountWithoutAuth(r, config, instanceCtrl, authCtrl)
		setupWebhooks(r, tenantCtrl, projectCtrl, ghAppCtrl)
		r.Group(func(r chi.Router) {
			r.Use(middlewareauthn.AttemptWeb(authenticator, instanceCtrl, authCtrl, config.Token.CookieName))
			setupRoutesV1WithAuth(r,
				appCtx, config,
				instanceCtrl, serverCtrl,
				userCtrl, tenantCtrl, projectCtrl,
				envCtrl, authCtrl,
				ghAppCtrl, gitPublicCtrl,
				appCtrl, varCtrl,
				deploymentCtrl, logsCtrl,
				volumeCtrl, templCtrl,
				favCtrl,
			)
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			render.NotFound(w, r)
		})
	})

	return r
}

// nolint: revive // it's the app context, it shouldn't be the first argument
func setupRoutesV1WithAuth(r chi.Router,
	appCtx context.Context,
	config *types.Config,
	instanceCtrl *instance.Controller,
	serverCtrl *server.Controller,
	userCtrl *user.Controller,
	tenantCtrl *tenant.Controller,
	projectCtrl *project.Controller,
	envCtrl *environment.Controller,
	authCtrl *auth.Controller,
	ghAppCtrl *githubapp.Controller,
	gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller,
	varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller,
	volumeCtrl *volume.Controller,
	templCtrl *template.Controller,
	favCtrl *favorite.Controller,
) {

	setupAccount(r, config, authCtrl, userCtrl, tenantCtrl)

	//TODO: multi tenant level routes goes here

	//Personal tenant routes
	setupInstance(r, instanceCtrl, serverCtrl, authCtrl)
	setupTenant(r, appCtx, tenantCtrl, projectCtrl, envCtrl, ghAppCtrl, gitPublicCtrl, appCtrl, varCtrl, deploymentCtrl, logsCtrl, volumeCtrl, templCtrl, favCtrl)
}

func setupWebhooks(r chi.Router, tenantCtrl *tenant.Controller, projectCtrl *project.Controller, ghAppCtrl *githubapp.Controller) {
	r.Route("/webhooks", func(r chi.Router) {
		r.Route("/source", func(r chi.Router) {
			r.Route("/github", func(r chi.Router) {
				r.Get("/redirect", webhookhandler.HandleGithubRedirect(tenantCtrl, projectCtrl, ghAppCtrl))
				r.Get("/install", webhookhandler.HandleGithubInstall(tenantCtrl, projectCtrl, ghAppCtrl))
			})
		})
	})
}

func setupInstance(r chi.Router, instanceCtrl *instance.Controller, serverCtrl *server.Controller, authCtrl *auth.Controller) {
	r.Route("/settings", func(r chi.Router) {
		r.Use(middlewarerestrict.ToSuperAdmin())
		r.Use(middlewarenav.PopulateNavItemKey("Instance Settings"))
		r.Get("/", handlerinstance.HandleGetSettings(instanceCtrl, serverCtrl))
		r.Patch("/fqdn", handlerinstance.HandlePatchFQDN(instanceCtrl, serverCtrl))
		r.Patch("/dns", handlerinstance.HandlePatchDNS(instanceCtrl, serverCtrl))
		r.Patch("/scripts", handlerinstance.HandlePatchScripts(instanceCtrl, serverCtrl))
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", handlerinstance.HandleGetAuth(instanceCtrl, authCtrl))
			r.Patch("/password", handlerinstance.HandlePatchPassword(instanceCtrl, authCtrl))
			r.Patch("/demo", handlerinstance.HandlePatchDemoUser(instanceCtrl, authCtrl))
			r.Patch("/github", handlerinstance.HandlePatchOauthProvider(instanceCtrl, authCtrl, enum.AuthProviderGithub))
		})
		r.Route("/registry", func(r chi.Router) {
			r.Get("/", handlerinstance.HandleGetRegistry(instanceCtrl))
			r.Patch("/", handlerinstance.HandlePatchRegistry(instanceCtrl))
		})

		setupServer(r, serverCtrl)
	})
}

func setupServer(r chi.Router, serverCtrl *server.Controller) {
	r.Route("/server", func(r chi.Router) {
		r.Use(middlewarerestrict.ToSuperAdmin())
		r.Get("/", handlerserver.HandleGet(serverCtrl))
		r.Patch("/", handlerserver.HandlePatchGeneral(serverCtrl))
		r.Patch("/network", handlerserver.HandlePatchNetwork(serverCtrl))
		r.Patch("/builder", handlerserver.HandlePatchBuilder(serverCtrl))
		r.Patch("/limits", handlerserver.HandlePatchLimits(serverCtrl))
		r.Get("/certificates", handlerserver.HandleListCertificates(serverCtrl))
	})
}

func setupTenant(r chi.Router,
	appCtx context.Context,
	tenantCtrl *tenant.Controller,
	projectCtrl *project.Controller,
	envCtrl *environment.Controller,
	ghAppCtrl *githubapp.Controller,
	gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller,
	varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller,
	volumeCtrl *volume.Controller,
	templCtrl *template.Controller,
	favCtrl *favorite.Controller,
) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { render.RedirectWithRefresh(w, "/team") })
		r.Route("/team", func(r chi.Router) {
			r.Use(middlewarenav.PopulateNavItemKey("Team"))
			r.Get("/", handlertenant.HandleList(tenantCtrl))
			r.Get(fmt.Sprintf("/nav/{%s}", request.PathParamSelectedUID), handlertenant.HandleListNavigation(tenantCtrl))
			r.Route("/new", func(r chi.Router) {
				r.Use(middlewarerestrict.ToSuperAdmin())
				r.Use(middlewarenav.PopulateNavItemKey("New Team"))
				r.Get("/", handlertenant.HandleNew())
				r.Post("/", handlertenant.HandleAdd(tenantCtrl))
			})
			r.Route(fmt.Sprintf("/{%s}", request.PathParamTenantUID), func(r chi.Router) {
				r.Use(middlewareinject.InjectTenant(tenantCtrl))
				r.Use(middlewarenav.PopulateNavTeam())
				r.Get("/", handlertenant.HandleGet(tenantCtrl, projectCtrl))
				r.Get("/favorites", handlerfavorite.HandleListFavorites(favCtrl))
				setupProject(r, appCtx, tenantCtrl, projectCtrl, envCtrl, ghAppCtrl, gitPublicCtrl, appCtrl, varCtrl, deploymentCtrl, logsCtrl, volumeCtrl, templCtrl, favCtrl)

				// Admin routes
				r.Route("/", func(r chi.Router) {
					r.Use(middlewarerestrict.ToTeamAdmin())
					r.Get("/settings", handlertenant.HandleGetSettings(tenantCtrl))
					r.Patch("/settings", handlertenant.HandlePatchGeneralSettings(tenantCtrl))
					r.Get("/restrictions", handlertenant.HandleGetRestrictions(tenantCtrl))
					r.Patch("/restrictions", handlertenant.HandlePatchRestrictions(tenantCtrl))
					r.Route("/members", func(r chi.Router) {
						r.Get("/", handlertenant.HandleListMembers(tenantCtrl))
						r.Post("/", handlertenant.HandleAddMember(tenantCtrl))
						r.Patch("/", handlertenant.HandlePatchMember(tenantCtrl))
						r.Delete("/", handlertenant.HandleDeleteMember(tenantCtrl))
					})
					r.Delete("/delete", handlertenant.HandleDeleteTeam(tenantCtrl))
				})
			})
		})
	})
}

func setupProject(r chi.Router,
	appCtx context.Context,
	tenantCtrl *tenant.Controller,
	projectCtrl *project.Controller,
	envCtrl *environment.Controller,
	ghAppCtrl *githubapp.Controller,
	gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller,
	varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller,
	volumeCtrl *volume.Controller,
	templCtrl *template.Controller,
	favCtrl *favorite.Controller,
) {
	r.Route("/project", func(r chi.Router) {
		r.Route("/new", func(r chi.Router) {
			r.Use(middlewarerestrict.ToTeamAdmin())
			r.Use(middlewarenav.PopulateNavItemKey("New Project"))
			r.Get("/", handlerproject.HandleNew())
			r.Post("/", handlerproject.HandleAdd(projectCtrl))
		})
		r.Get(fmt.Sprintf("/nav/{%s}", request.PathParamSelectedUID), handlerproject.HandleListNavigation(projectCtrl))
		r.Route(fmt.Sprintf("/{%s}", request.PathParamProjectUID), func(r chi.Router) {
			r.Use(middlewareinject.InjectProject(projectCtrl))
			r.Use(middlewarerestrict.ToProjectRole())
			r.Use(middlewarenav.PopulateNavProject())
			r.Route("/", func(r chi.Router) {
				r.Use(middlewarerestrict.ModificationToProjectOwner()) // only owners can modify, others can view
				r.Get("/", handlerproject.HandleGet(projectCtrl, envCtrl, appCtrl))
				r.Get("/overview", handlerproject.HandleGet(projectCtrl, envCtrl, appCtrl))
				r.Get("/settings", handlerproject.HandleGetSettingsGeneral(projectCtrl))
				r.Patch("/settings", handlerproject.HandleUpdateSettingsGeneral(projectCtrl))
				r.Delete("/delete", handlerproject.HandleDelete(projectCtrl))
			})
			r.Get("/events", handlerproject.HandleEvents(appCtx, projectCtrl))
			setupEnvionment(r, appCtx, envCtrl, ghAppCtrl, gitPublicCtrl, appCtrl, varCtrl, deploymentCtrl, logsCtrl, volumeCtrl, templCtrl, favCtrl)
			setupProjectSource(r, ghAppCtrl, gitPublicCtrl)

			// Admin/Owner routes
			r.Route("/members", func(r chi.Router) {
				r.Use(middlewarerestrict.ToProjectOwner())
				r.Get("/", handlerproject.HandleListMembers(projectCtrl))
				r.Post("/", handlerproject.HandleAddMember(projectCtrl))
				r.Patch("/", handlerproject.HandleUpdateMember(projectCtrl))
				r.Delete("/", handlerproject.HandleDeleteMember(projectCtrl))
				r.Get("/list-nonmembers", handlerproject.HandleListAllMembers(tenantCtrl))
			})
		})
	})
}

func setupEnvionment(r chi.Router,
	appCtx context.Context,
	envCtrl *environment.Controller, ghAppCtrl *githubapp.Controller, gitPublicCtrl *gitpublic.Controller,
	appCtrl *application.Controller, varCtrl *variable.Controller,
	deploymentCtrl *deployment.Controller,
	logsCtrl *logs.Controller, volumeCtrl *volume.Controller,
	templCtrl *template.Controller, favCtrl *favorite.Controller,
) {
	r.Route("/environment", func(r chi.Router) {
		r.Get("/", handlerenvironment.HandleList(envCtrl))
		r.Post("/create", handlerenvironment.HandleAdd(envCtrl))
		r.Get("/nav", handlerenvironment.HandleListNavigation(envCtrl))
		r.Route(fmt.Sprintf("/{%s}", request.PathParamEnvironmentUID), func(r chi.Router) {
			r.Use(middlewareinject.InjectEnvironment(envCtrl))
			r.Route("/", func(r chi.Router) {
				r.Use(middlewarerestrict.ModificationToProjectOwner()) // only owners can modify, others can view
				r.Patch("/settings", handlerenvironment.HandleUpdate(envCtrl))
				r.Route("/volumes", func(r chi.Router) {
					r.Get("/", handlervolume.HandleListUnattached(volumeCtrl))
					r.Route(fmt.Sprintf("/{%s}", request.PathParamVolumeUID), func(r chi.Router) {
						r.Use(middlewareinject.InjectVolume(volumeCtrl))
						r.Delete("/", handlervolume.HandleDelete(volumeCtrl))
					})
				})
				r.Delete("/delete", handlerenvironment.HandleDelete(envCtrl))
			})
			setupApplication(r, appCtx, envCtrl, appCtrl, varCtrl, ghAppCtrl, gitPublicCtrl, deploymentCtrl, logsCtrl, volumeCtrl, templCtrl, favCtrl)
		})
	})
}

func setupApplication(r chi.Router, appCtx context.Context, envCtrl *environment.Controller, appCtrl *application.Controller, varCtrl *variable.Controller, ghAppCtrl *githubapp.Controller, gitPublicCtrl *gitpublic.Controller, deploymentCtrl *deployment.Controller, logsCtrl *logs.Controller, volumeCtrl *volume.Controller, templCtrl *template.Controller, favCtrl *favorite.Controller) {
	r.Route("/application", func(r chi.Router) {
		r.Get("/", handlerapplication.HandleList(envCtrl, appCtrl))
		r.Get("/nav", handlerapplication.HandleListNavigation(appCtrl))
		r.Route(fmt.Sprintf("/{%s}", request.PathParamApplicationUID), func(r chi.Router) {
			//Inject application here
			r.Use(middlewareinject.InjectApplication(appCtrl))
			r.Get("/", handlerapplication.HandleListDeployments(appCtrl, deploymentCtrl))
			r.Patch("/name", handlerapplication.HandleUpdateName(appCtrl))
			r.Patch("/icon", handlerapplication.HandleUpdateIcon(appCtrl))
			r.Get("/deployments", handlerapplication.HandleListDeployments(appCtrl, deploymentCtrl))
			r.Get("/settings", handlerapplication.HandleGetSettings(appCtrl, ghAppCtrl))
			r.Patch("/settings", handlerapplication.HandleUpdateSettings(appCtrl, ghAppCtrl))
			r.Get("/logs", handlerapplication.HandleGetLogs(appCtrl))
			r.Get("/logs/stream", handlerapplication.HandleTailLogs(appCtx, appCtrl))
			r.Get("/terminal", handlerapplication.HandleGetTerminal())
			r.Get(fmt.Sprintf("/metrics/{%s}", request.PathParamMetricsSpan), handlerapplication.HandleGetMetrics(appCtrl))
			r.Route("/favorite", func(r chi.Router) {
				r.Get("/", handlerfavorite.HandleGetFavorite(favCtrl))
				r.Get("/add", handlerfavorite.HandleAddFavorite(favCtrl))
				r.Get("/remove", handlerfavorite.HandleDeleteFavorite(favCtrl))
			})
			r.Route("/network", func(r chi.Router) {
				r.Route("/http", func(r chi.Router) {
					r.Post("/generate", handlerapplication.HandleGenerateDomain(appCtrl, ghAppCtrl))
					r.Post("/", handlerapplication.HandleUpdateDomain(appCtrl, ghAppCtrl))
					r.Delete("/", handlerapplication.HandleDeleteDomain(appCtrl, ghAppCtrl))
				})
				r.Route("/tcp", func(r chi.Router) {
					r.Post("/", handlerapplication.HandleAddTCPProxy(appCtrl, ghAppCtrl))
					r.Delete("/", handlerapplication.HandleDeleteTCPProxy(appCtrl, ghAppCtrl))
				})
				r.Post("/private", handlerapplication.HandleUpdatePrivateDomain(appCtrl, ghAppCtrl))
			})
			r.Route("/variables", func(r chi.Router) {
				r.Get("/", handlervariable.HandleList(varCtrl))
				r.Post("/", handlervariable.HandlePost(appCtrl, varCtrl))
				r.Route(fmt.Sprintf("/{%s}", request.PathParamVariableUID), func(r chi.Router) {
					r.Patch("/", handlervariable.HandlePatch(appCtrl, varCtrl))
					r.Patch("/generate", handlervariable.HandleGenerate(appCtrl, varCtrl))
					r.Delete("/", handlervariable.HandleDelete(appCtrl, varCtrl))
				})
			})
			r.Route("/volumes", func(r chi.Router) {
				r.Get("/", handlervolume.HandleListVolume(volumeCtrl))
				r.Post("/create", handlervolume.HandleCreate(appCtrl, volumeCtrl))
				r.Get("/unattached", handlervolume.HandleListAttachable(volumeCtrl))
				r.Route(fmt.Sprintf("/{%s}", request.PathParamVolumeUID), func(r chi.Router) {
					r.Use(middlewareinject.InjectVolume(volumeCtrl))
					r.Patch("/", handlervolume.HandleUpdateAttached(appCtrl, volumeCtrl))
					r.Patch("/detach", handlervolume.HandleUpdateDetach(appCtrl, volumeCtrl))
				})
			})
			r.Patch("/deploy", handlerapplication.HandleDeploy(appCtrl))
			r.Patch("/redeploy", handlerapplication.HandleRedeploy(appCtrl))
			r.Get("/delete", handlerapplication.HandleDeleteView())
			r.Delete("/delete", handlerapplication.HandleDeleteApplication(appCtrl))
			setupDeployment(r, appCtx, appCtrl, deploymentCtrl, logsCtrl)
		})
		setupApplicationCreate(r, appCtrl, ghAppCtrl, gitPublicCtrl, templCtrl)
	})
}

func setupApplicationCreate(r chi.Router, appCtrl *application.Controller, ghAppCtrl *githubapp.Controller, gitPublicCtrl *gitpublic.Controller, templCtrl *template.Controller) {
	r.Route("/new", func(r chi.Router) {
		r.Get("/", handlercreate.HandleListGitOptions(dto.SourceCategoryGit))
		r.Route("/git", func(r chi.Router) {
			r.Get("/", handlercreate.HandleListGitOptions(dto.SourceCategoryGit))
			r.Route("/github", func(r chi.Router) {
				r.Get("/", handlercreate.HandleListGithubApps(ghAppCtrl))
				r.Route(fmt.Sprintf("/{%s}", request.PathParamSourceUID), func(r chi.Router) {
					r.Use(middlewareinject.InjectGithubAppSource(ghAppCtrl))
					r.Get("/", handlercreate.HandleGetGithubView(appCtrl))
					r.Post("/", handlercreate.HandleCreateGithub(appCtrl))
				})
			})
			r.Route("/git-public", func(r chi.Router) {
				r.Get("/", handlercreate.HandleGetGitPublicView())
				r.Post("/load-repo", handlercreate.HandleLoadGitPublicRepo(appCtrl))
				r.Post("/", handlercreate.HandleCreateGitPublic(appCtrl))
			})
		})
		r.Route("/registry", func(r chi.Router) {
			r.Get("/", handlercreate.HandleGetRegistryView())
			r.Post("/", handlercreate.HandleCreateWithRegistry(appCtrl))
		})
		r.Route("/database", func(r chi.Router) {
			r.Get("/", handlercreate.HandleGetDatabaseView(templCtrl))
			r.Route(fmt.Sprintf("/{%s}", request.PathParamTemplateID), func(r chi.Router) {
				r.Post("/", handlercreate.HandleTemplateCreate(templCtrl))
				r.Get("/preview", handlercreate.HandleTemplatePreview(templCtrl, "database"))
			})
		})
		r.Route("/oneclick", func(r chi.Router) {
			r.Get("/", handlercreate.HandleGetOnelickTemplate(templCtrl))
			r.Route(fmt.Sprintf("/{%s}", request.PathParamTemplateID), func(r chi.Router) {
				r.Post("/", handlercreate.HandleTemplateCreate(templCtrl))
				r.Get("/preview", handlercreate.HandleTemplatePreview(templCtrl, "oneclick"))
			})
		})
	})
}

func setupDeployment(r chi.Router, appCtx context.Context, appCtrl *application.Controller, deploymentCtrl *deployment.Controller, logsCtrl *logs.Controller) {
	r.Route("/deployment", func(r chi.Router) {
		r.Route(fmt.Sprintf("/{%s}", request.PathParamDeploymentUID), func(r chi.Router) {
			//Inject deployment here
			r.Use(middlewareinject.InjectDeployment(deploymentCtrl))
			r.Get("/", handlerdeployment.HandleGetDeployment(deploymentCtrl))
			r.Get("/logs", handlerLogs.HandleGetLogs(logsCtrl))
			r.Get("/logs/stream", handlerLogs.HandleTailLogs(appCtx, logsCtrl))
		})
	})
}

func setupProjectSource(r chi.Router, ghAppCtrl *githubapp.Controller, gitPublicCtrl *gitpublic.Controller) {
	r.Route("/source", func(r chi.Router) {
		r.Get("/", handlersource.HandleListConfigurableSource())
		r.Route("/git-public", func(r chi.Router) {
			r.Get("/list-branches", handlersource.HandleListGitpublicBranches(gitPublicCtrl))
		})
		r.Route("/github", func(r chi.Router) {
			r.Get("/", handlersource.HandleListGithubApps(ghAppCtrl))
			r.Post("/", handlersource.HandleAddGithubApp(ghAppCtrl))
			r.Route(fmt.Sprintf("/{%s}", request.PathParamSourceUID), func(r chi.Router) {
				r.Use(middlewareinject.InjectGithubAppSource(ghAppCtrl))
				r.Route("/", func(r chi.Router) {
					r.Use(middlewarerestrict.ToProjectOwner())
					r.Get("/", handlersource.HandleGetGithubApp(ghAppCtrl))
					r.Delete("/", handlersource.HandleDeleteGithubApp(ghAppCtrl))
				})
				r.Get("/list-repos", handlersource.HandleListGithubRepos(ghAppCtrl))
				r.Get("/list-branches", handlersource.HandleListGithubBranches(ghAppCtrl))
			})
		})
	})
}

func setupAccountWithoutAuth(r chi.Router, config *types.Config, instanceCtrl *instance.Controller, authCtrl *auth.Controller) {
	cookieName := config.Token.CookieName
	r.Get("/login", authhandler.HandleGet(authCtrl, instanceCtrl, cookieName))
	r.Post("/login", authhandler.HandleLoginPost(authCtrl, cookieName))
	r.Get("/register", authhandler.HandleGetRegister(instanceCtrl))
	r.Post("/register", authhandler.HandleRegister(instanceCtrl, authCtrl, cookieName))
	r.Route("/auth", func(r chi.Router) {
		r.Route(fmt.Sprintf("/{%s}", request.PathParamAuthProvider), func(r chi.Router) {
			r.Use(middlewareinject.InjectAuthProvider(authCtrl))
			r.Get("/redirect", authhandler.HandleRedirect(authCtrl))
			r.Get("/callback", authhandler.HandleCallback(authCtrl, cookieName))
		})
	})
}

func setupAccount(r chi.Router, config *types.Config, authCtrl *auth.Controller, userCtrl *user.Controller, tenantCtrl *tenant.Controller) {
	cookieName := config.Token.CookieName
	r.Get("/logout", accounthandler.HandleLogout(authCtrl, cookieName))
	r.Route("/account", func(r chi.Router) {
		r.Use(middlewarenav.PopulateNavItemKey("Account"))
		r.Get("/", accounthandler.HandleGetProfile(userCtrl))
		r.Patch("/", accounthandler.HandlePatchProfile(userCtrl))
		r.Get("/session", accounthandler.HandleGetSession(userCtrl))
		r.Get("/delete", accounthandler.HandleGetDelete(userCtrl))
	})
}

func staticDev() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("./app/web/public")))
}

func staticProd() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(public.AssetsFS))
}

func enableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=1440, public")
		next.ServeHTTP(w, r)
	})
}

func disableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
