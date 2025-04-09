import effects.{fetch_messages}
import gleam/dict
import gleam/list
import gleam/option
import lustre/effect
import model.{type Model, Model}
import types.{type Msg, Behaviours}

const message_count_interval_secs = 5

pub fn update(model: Model, msg: Msg) -> #(Model, effect.Effect(Msg)) {
  case msg {
    types.ConfigFetched(res) ->
      case res {
        Error(e) -> #(Model(..model, http_error: option.Some(e)), effect.none())
        Ok(c) -> #(Model(..model, config: option.Some(c)), effect.none())
      }
    types.BehavioursFetched(res) ->
      case res {
        Error(_) -> #(model, effect.none())
        Ok(b) ->
          case b.show_live_count {
            False -> #(Model(..model, behaviours: b), effect.none())
            True -> #(
              Model(..model, behaviours: b),
              effect.batch([
                effects.fetch_message_count(),
                effects.schedule_next_tick(message_count_interval_secs),
              ]),
            )
          }
      }
    types.FetchMessages(num) ->
      case num {
        1 -> #(
          Model(..model, fetching: True, http_error: option.None),
          fetch_messages(num, model.behaviours.delete_messages),
        )
        _ -> #(
          Model(..model, fetching: True, http_error: option.None),
          fetch_messages(num, model.behaviours.delete_messages),
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
      Model(
        ..model,
        behaviours: Behaviours(..model.behaviours, select_on_hover: selected),
      ),
      effect.none(),
    )
    types.DeleteSettingsChanged(selected) -> #(
      Model(
        ..model,
        behaviours: Behaviours(..model.behaviours, delete_messages: selected),
      ),
      effect.none(),
    )
    types.ShowLiveCountChanged(selected) ->
      case selected {
        False -> #(
          Model(
            ..model,
            message_count: option.None,
            behaviours: Behaviours(
              ..model.behaviours,
              show_live_count: selected,
            ),
          ),
          effect.none(),
        )
        True -> #(
          Model(
            ..model,
            behaviours: Behaviours(
              ..model.behaviours,
              show_live_count: selected,
            ),
          ),
          effect.batch([
            effects.fetch_message_count(),
            effects.schedule_next_tick(message_count_interval_secs),
          ]),
        )
      }
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
    types.MessageCountFetched(res) ->
      case res {
        Error(_) -> #(Model(..model, message_count: option.None), effect.none())
        Ok(c) -> #(
          Model(..model, message_count: option.Some(c.count)),
          effect.none(),
        )
      }
    types.Tick ->
      case model.behaviours.show_live_count {
        False -> #(model, effect.none())
        True -> #(
          model,
          effect.batch([
            effects.fetch_message_count(),
            effects.schedule_next_tick(message_count_interval_secs),
          ]),
        )
      }
  }
}
