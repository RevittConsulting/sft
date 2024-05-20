Checks toggles via a call to CheckFeatureIsEnabled which is in the service.

So in order to use the library, it needs to be initialised with an sftService added to the
dependency array, which can then be used throughout the appliction, via the check function.

This is the code I've been using to carry out the check:

enabled, err := h.sft.CheckFeatureIsEnabled(r.Context(), "Add task")
if err != nil {
    log.Println("error checking feature: ", err)
}
if enabled.Enabled != true { 
    return
}



Here's how I've set it up in CHD:

in server/setup.go

Adding to the deps:
// toggles
sftDb := sft.NewDb(s.Pool)
sftService := sft.NewService(sftDb, context.Background(), s.Pool)
deps.sft = sftService

Setting up the routes:
// toggles stuff
sftConfig := &sft.Config{
Buildpath: "/Users/maxbb/github/revitt/sft/web/dashboard/dist",
Port:      "6969",
}
	s.Router.Route("/api/sft/v1", func(r chi.Router) {
		sft.NewHandler(r, s.Deps.sft, sftConfig)
	})


