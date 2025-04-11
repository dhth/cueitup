import gleam/dynamic/decode
import gleam/option
import lustre_http

pub type Config {
  Config(
    profile_name: String,
    queue_url: String,
    aws_config_source: String,
    context_key: option.Option(String),
    subset_key: option.Option(String),
  )
}

pub fn config_decoder() -> decode.Decoder(Config) {
  use profile_name <- decode.field("profile_name", decode.string)
  use queue_url <- decode.field("queue_url", decode.string)
  use aws_config_source <- decode.field("aws_config_source", decode.string)
  use context_key <- decode.field("context_key", decode.optional(decode.string))
  use subset_key <- decode.field("subset_key", decode.optional(decode.string))
  decode.success(Config(
    profile_name:,
    queue_url:,
    aws_config_source:,
    context_key:,
    subset_key:,
  ))
}

pub type Behaviours {
  Behaviours(
    delete_messages: Bool,
    select_on_hover: Bool,
    show_message_count: Bool,
  )
}

pub fn default_behaviours() -> Behaviours {
  Behaviours(
    delete_messages: True,
    select_on_hover: False,
    show_message_count: False,
  )
}

pub fn behaviours_decoder() -> decode.Decoder(Behaviours) {
  use delete_messages <- decode.field("delete_messages", decode.bool)
  use select_on_hover <- decode.field("select_on_hover", decode.bool)
  use show_message_count <- decode.field("show_message_count", decode.bool)
  decode.success(Behaviours(
    delete_messages:,
    select_on_hover:,
    show_message_count:,
  ))
}

pub type MessageOffset =
  Int

pub type Message {
  Message(
    id: String,
    body: String,
    context_key: option.Option(String),
    context_value: option.Option(String),
    error: option.Option(String),
  )
}

pub fn message_details_decoder() -> decode.Decoder(Message) {
  use id <- decode.field("id", decode.string)
  use body <- decode.field("body", decode.string)
  use context_key <- decode.field("context_key", decode.optional(decode.string))
  use context_value <- decode.field(
    "context_value",
    decode.optional(decode.string),
  )
  use error <- decode.field("error", decode.optional(decode.string))
  decode.success(Message(id:, body:, context_key:, context_value:, error:))
}

pub type MessageCount {
  MessageCount(count: Int)
}

pub fn message_count_decoder() -> decode.Decoder(MessageCount) {
  use count <- decode.field("count", decode.int)
  decode.success(MessageCount(count:))
}

pub type Msg {
  ConfigFetched(Result(Config, lustre_http.HttpError))
  BehavioursFetched(Result(Behaviours, lustre_http.HttpError))
  FetchMessages(Int)
  ClearMessages
  HoverSettingsChanged(Bool)
  DeleteSettingsChanged(Bool)
  ShowMessageCountChanged(Bool)
  MessageChosen(Int)
  MessagesFetched(Result(List(Message), lustre_http.HttpError))
  MessageCountFetched(Result(MessageCount, lustre_http.HttpError))
  GoToStart
  GoToEnd
  Tick
}

pub fn dummy_message() -> List(Message) {
  let id = "20693f56-b784-4594-b79e-38c6d1756035"
  let body =
    "
{
  \"aggregateId\": \"00000000-0000-0000-0000-000000012363\",
  \"dateTime\": \"2025-04-08T16:39:10.561+0000\",
  \"hash\": \"xodq\",
  \"longUrl\": \"https://github.com/rusqlite/rusqlite\",
  \"metaData\": {
    \"contentId\": \"b44f7d51-662f-4c8f-971a-081f6da0758b\",
    \"redirectMode\": \"DIRECT\",
    \"referrerId\": \"02b63ef8-e801-4995-a2a5-2dc128e747f4\",
    \"referrerUserId\": \"02b63ef8-e801-4995-a2a5-2dc128e747f4\",
    \"tenantId\": \"902daed1-911a-4641-9a11-00631e313274\"
  },
  \"sequenceNr\": 1,
  \"traceId\": \"d42e9ebc-8abb-49ea-b365-4e9bd7a6b0b0\",
  \"type\": \"eu.firstbird.hummingbird.url.event.UrlShortened\",
  \"version\": 4
}"
  [
    Message(
      id:,
      body:,
      context_key: option.None,
      context_value: option.None,
      error: option.None,
    ),
  ]
}
