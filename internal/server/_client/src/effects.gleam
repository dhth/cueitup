import gleam/dynamic/decode
import gleam/int
import lustre/effect
import lustre_http
import plinth/browser/window
import plinth/javascript/global
import types.{
  behaviours_decoder, config_decoder, message_count_decoder,
  message_details_decoder,
}

const dev = True

fn base_url() -> String {
  case dev {
    False -> window.location()
    True -> "http://127.0.0.1:8500/"
  }
}

pub fn fetch_config() -> effect.Effect(types.Msg) {
  let expect = lustre_http.expect_json(config_decoder(), types.ConfigFetched)

  lustre_http.get(base_url() <> "api/config", expect)
}

pub fn fetch_behaviours() -> effect.Effect(types.Msg) {
  let expect =
    lustre_http.expect_json(behaviours_decoder(), types.BehavioursFetched)

  lustre_http.get(base_url() <> "api/behaviours", expect)
}

pub fn fetch_message_count() -> effect.Effect(types.Msg) {
  let expect =
    lustre_http.expect_json(message_count_decoder(), types.MessageCountFetched)
  lustre_http.get(base_url() <> "api/message-count", expect)
}

pub fn fetch_messages(num: Int, delete: Bool) -> effect.Effect(types.Msg) {
  let expect =
    lustre_http.expect_json(
      decode.list(message_details_decoder()),
      types.MessagesFetched,
    )

  let delete_query_param = case delete {
    False -> "false"
    True -> "true"
  }

  lustre_http.get(
    base_url()
      <> "api/fetch?num="
      <> num |> int.to_string
      <> "&delete="
      <> delete_query_param,
    expect,
  )
}

pub fn schedule_next_tick(delay_seconds: Int) -> effect.Effect(types.Msg) {
  effect.from(fn(dispatch) {
    global.set_timeout(delay_seconds * 1000, fn() { dispatch(types.Tick) })
    Nil
  })
}
