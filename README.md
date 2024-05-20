<div align="center" style="padding-top: 30px; padding-bottom: 30px">
    <h1>SFT</h1>
</div>

Simple Feature Toggles. A simple library for toggling features within your application.

***

## Usage

Create a new instance of sft
- `NewDb()` takes a pgxpool
- `NewService()` takes a NewBd() (which implements ISimpleFeatureToggleDb), context.Context, and a pgxpool.Pool.

Manage toggles via the dashboard, which is accessed at: `......:6969/dashboard`

Once a toggle is created, insert the toggle check into the code of the relevant feature, using `CheckFeatureIsEnabled()`, which takes a context.Context, and a feature name (string).

Example toggle checking logic:

```
    enabled, err := h.sft.CheckFeatureIsEnabled(r.Context(), "fix issue")
	if err != nil {
		log.Println("error checking feature: ", err)
	}
	if enabled.Enabled == false {
		log.Println("feature is currently disabled")
		return
	}
```