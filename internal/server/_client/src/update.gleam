import effects.{fetch_messages}
import gleam/dict
import gleam/list
import gleam/option
import lustre/effect
import model.{type Model, Model}
import types.{type Msg}

pub fn update(model: Model, msg: Msg) -> #(Model, effect.Effect(Msg)) {
  case msg {
    types.ConfigFetched(res) ->
      case res {
        Error(e) -> #(Model(..model, http_error: option.Some(e)), effect.none())
        Ok(c) -> #(Model(..model, config: option.Some(c)), effect.none())
      }
    types.FetchMessages(num) ->
      case num {
        1 -> #(
          Model(..model, fetching: True, http_error: option.None),
          fetch_messages(num),
        )
        _ -> #(
          Model(..model, fetching: True, http_error: option.None),
          fetch_messages(num),
        )
      }
    types.ClearMessages -> #(
      Model(
        ..model,
        messages: [],
        messages_cache: dict.new(),
        current_message: option.None,
        http_error: option.None,
      ),
      effect.none(),
    )
    types.HoverSettingsChanged(selected) -> #(
      Model(..model, select_on_hover: selected),
      effect.none(),
    )
    types.GoToEnd -> #(model, effect.none())
    types.GoToStart -> #(model, effect.none())
    types.MessageChosen(index) -> {
      let maybe_message = model.messages_cache |> dict.get(index)
      case maybe_message {
        Error(_) -> #(model, effect.none())
        Ok(msg) -> #(
          Model(..model, current_message: option.Some(#(index, msg))),
          effect.none(),
        )
      }
    }
    types.MessagesFetched(result) ->
      case result {
        Error(e) -> #(
          Model(..model, fetching: False, http_error: option.Some(e)),
          effect.none(),
        )
        Ok(messages) -> {
          let updated_messages = model.messages |> list.append(messages)
          let messages_cache =
            updated_messages
            |> list.index_map(fn(m, i) { #(i, m) })
            |> dict.from_list
          #(
            Model(
              ..model,
              fetching: False,
              messages: updated_messages,
              messages_cache: messages_cache,
            ),
            effect.none(),
          )
        }
      }
  }
}
