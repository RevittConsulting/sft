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

Find an example in fix issues.