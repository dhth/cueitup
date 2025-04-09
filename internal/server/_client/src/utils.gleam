import gleam/dynamic/decode
import gleam/int
import gleam/json
import gleam/list
import gleam/string
import lustre_http

pub fn http_error_to_string(error: lustre_http.HttpError) -> String {
  case error {
    lustre_http.BadUrl(u) -> "bad url:" <> u
    lustre_http.InternalServerError(e) -> "internal server error: " <> e
    lustre_http.JsonError(e) ->
      case e {
        json.UnableToDecode(de) ->
          de
          |> list.map(fn(err) {
            case err {
              decode.DecodeError(exp, found, _) ->
                "couldn't decode JSON; expected: "
                <> exp
                <> ", found: "
                <> found
            }
          })
          |> string.join(", ")
        json.UnexpectedByte(_) -> "unexpected byte"
        json.UnexpectedEndOfInput -> "unexpected end of input"
        json.UnexpectedFormat(_) -> "unexpected format"
        json.UnexpectedSequence(_) -> "unexpected sequence"
      }
    lustre_http.NetworkError -> "network error"
    lustre_http.NotFound -> "not found"
    lustre_http.OtherError(code, body) ->
      "non success HTTP response; status: "
      <> int.to_string(code)
      <> ", body: "
      <> body
    lustre_http.Unauthorized -> "unauthorized"
  }
}
