type StringDate = string
type Title = string
type NullableTitle = Title | null
type Tag = string
type Time = StringDate
type NullableTime = Time | null
type URL = string
type NullableURL = URL | null
export type Article = {
    title: NullableTitle
    tags: Tag[]
    time: NullableTime
    url: NullableURL
}
