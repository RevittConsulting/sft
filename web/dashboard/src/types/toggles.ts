export type Toggle = {
    id: string
    feature_name: string
    toggle_meta: Record<string, string>
    enabled: boolean
}

export type ToggleDto = {
    feature_name: string
    toggle_meta: Record<string, string>
    enabled: boolean
}