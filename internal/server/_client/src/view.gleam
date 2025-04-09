import gleam/int
import gleam/list
import gleam/option
import gleam/string
import lustre/attribute
import lustre/element
import lustre/element/html
import lustre/event
import model.{type Model}
import types.{type Config, type MessageDetails, type Msg}
import utils.{http_error_to_string}

const profile_name_max_width = 60

pub fn view(model: Model) -> element.Element(Msg) {
  html.div([attribute.class("bg-[#282828] text-[#ebdbb2] mt-4 mx-4")], [
    html.div([], [
      html.div([], [
        messages_section(model),
        controls_section(model),
        error_section(model),
      ]),
    ]),
  ])
}

fn messages_section(model: Model) -> element.Element(Msg) {
  let height_class = case model.http_error {
    option.None -> "h-[calc(100vh-4.3rem)]"
    option.Some(_) -> "h-[calc(100vh-9rem)]"
  }
  case model.messages {
    [] -> messages_section_empty(height_class)
    [_, ..] -> messages_section_with_messages(model, height_class)
  }
}

fn messages_section_empty(height_class: String) -> element.Element(Msg) {
  html.div(
    [
      attribute.class(
        "mt-4 "
        <> height_class
        <> " flex border-2 border-[#928374] border-opacity-20 items-center flex justify-center overflow-auto",
      ),
    ],
    [
      html.pre([attribute.class("text-[#928374]")], [
        element.text(
          "
                                                                                                                            
                                                                                                                            
                                                           iiii          tttt                                               
                                                          i::::i      ttt:::t                                               
                                                           iiii       t:::::t                                               
                                                                      t:::::t                                               
    ccccccccccccccccuuuuuu    uuuuuu      eeeeeeeeeeee   iiiiiiittttttt:::::ttttttt    uuuuuu    uuuuuu ppppp   ppppppppp   
  cc:::::::::::::::cu::::u    u::::u    ee::::::::::::ee i:::::it:::::::::::::::::t    u::::u    u::::u p::::ppp:::::::::p  
 c:::::::::::::::::cu::::u    u::::u   e::::::eeeee:::::eei::::it:::::::::::::::::t    u::::u    u::::u p:::::::::::::::::p 
c:::::::cccccc:::::cu::::u    u::::u  e::::::e     e:::::ei::::itttttt:::::::tttttt    u::::u    u::::u pp::::::ppppp::::::p
c::::::c     cccccccu::::u    u::::u  e:::::::eeeee::::::ei::::i      t:::::t          u::::u    u::::u  p:::::p     p:::::p
c:::::c             u::::u    u::::u  e:::::::::::::::::e i::::i      t:::::t          u::::u    u::::u  p:::::p     p:::::p
c:::::c             u::::u    u::::u  e::::::eeeeeeeeeee  i::::i      t:::::t          u::::u    u::::u  p:::::p     p:::::p
c::::::c     cccccccu:::::uuuu:::::u  e:::::::e           i::::i      t:::::t    ttttttu:::::uuuu:::::u  p:::::p    p::::::p
c:::::::cccccc:::::cu:::::::::::::::uue::::::::e         i::::::i     t::::::tttt:::::tu:::::::::::::::uup:::::ppppp:::::::p
 c:::::::::::::::::c u:::::::::::::::u e::::::::eeeeeeee i::::::i     tt::::::::::::::t u:::::::::::::::up::::::::::::::::p 
  cc:::::::::::::::c  uu::::::::uu:::u  ee:::::::::::::e i::::::i       tt:::::::::::tt  uu::::::::uu:::up::::::::::::::pp  
    cccccccccccccccc    uuuuuuuu  uuuu    eeeeeeeeeeeeee iiiiiiii         ttttttttttt      uuuuuuuu  uuuup::::::pppppppp    
                                                                                                         p:::::p            
                                                                                                         p:::::p            
                                                                                                        p:::::::p           
                                                                                                        p:::::::p           
                                                                                                        p:::::::p           
                                                                                                        ppppppppp

                    cueitup lets you inspect messages in an AWS SQS queue in a simple and deliberate manner

                                Click on the buttons below to start fetching messages
",
        ),
      ]),
    ],
  )
}

fn messages_section_with_messages(
  model: Model,
  height_class: String,
) -> element.Element(Msg) {
  let current_index =
    model.current_message
    |> option.map(fn(a) {
      case a {
        #(i, _) -> i
      }
    })

  html.div(
    [
      attribute.class(
        "mt-4 "
        <> height_class
        <> " flex border-2 border-[#928374] border-opacity-20",
      ),
    ],
    [
      html.div([attribute.class("w-2/5 overflow-auto")], [
        html.div([attribute.class("p-4")], [
          html.h2([attribute.class("text-[#d3869b] text-xl font-bold mb-4")], [
            html.text("Messages"),
          ]),
          html.div(
            [],
            model.messages
              |> list.index_map(fn(m, i) {
                message_list_item(
                  m,
                  i,
                  current_index,
                  model.behaviours.select_on_hover,
                )
              }),
          ),
        ]),
      ]),
      message_details_pane(model),
    ],
  )
}

fn message_list_item(
  message: MessageDetails,
  index: Int,
  current_index: option.Option(Int),
  select_on_hover: Bool,
) -> element.Element(Msg) {
  let border_class = case current_index {
    option.Some(i) if i == index -> " text-[#d3869b] border-l-[#d3869b]"
    _ -> " text-[#d5c4a1] border-l-[#282828]"
  }
  let event_handler = case select_on_hover {
    False -> event.on_click(types.MessageChosen(index))
    True -> event.on_mouse_over(types.MessageChosen(index))
  }

  html.div(
    [
      attribute.class(
        "py-2 px-4 border-l-2 hover:border-l-[#83a598]"
        <> " hover:text-[#83a598] hover:border-l-2 cursor-pointer transition duration-100"
        <> " ease-in-out"
        <> border_class,
      ),
      event_handler,
    ],
    [
      html.p([attribute.class("text-base font-semibold")], [
        html.text(message.id),
      ]),
      case message.context_key, message.context_value {
        option.Some(k), option.Some(v) ->
          html.div([attribute.class("flex space-x-2 text-sm")], [
            html.p([], [html.text(k <> ": " <> v)]),
          ])
        _, _ -> element.none()
      },
    ],
  )
}

fn message_details_pane(model: Model) -> element.Element(Msg) {
  let message_details = case model.current_message {
    option.None ->
      html.p([attribute.class("text-[#928374]")], [
        html.text(
          case model.behaviours.select_on_hover {
            True -> "Hover on"
            False -> "Select"
          }
          <> " an entry in the left pane to view details here.",
        ),
      ])
    option.Some(#(_, msg)) ->
      html.div([], [
        html.pre([attribute.class("text-[#d5c4a1] text-base mb-4")], [
          html.text(msg.body),
        ]),
      ])
  }

  html.div([attribute.class("w-3/5 p-6 overflow-auto")], [
    html.h2([attribute.class("text-[#d3869b] text-xl font-bold mb-4")], [
      html.text("Details"),
    ]),
    message_details,
  ])
}

fn controls_section(model: Model) -> element.Element(Msg) {
  case model.config {
    option.Some(c) -> controls_div_with_config(model, c)
    option.None -> controls_div_when_no_config()
  }
}

fn controls_div_when_no_config() -> element.Element(Msg) {
  html.div([attribute.class("flex items-center space-x-2 mt-4")], [
    html.button(
      [
        attribute.class(
          "font-bold px-4 py-1 bg-[#d3869b] text-[#282828] hover:bg-[#d3869b]",
        ),
        attribute.disabled(True),
      ],
      [
        html.a(
          [
            attribute.href("https://github.com/dhth/cueitup"),
            attribute.target("_blank"),
          ],
          [element.text("cueitup")],
        ),
      ],
    ),
  ])
}

fn controls_div_with_config(
  model: Model,
  config: Config,
) -> element.Element(Msg) {
  html.div([attribute.class("flex items-center space-x-2 mt-4")], [
    html.button(
      [
        attribute.class(
          "font-bold px-4 py-1 bg-[#d3869b] text-[#282828] hover:bg-[#d3869b]",
        ),
        attribute.disabled(True),
      ],
      [
        html.a(
          [
            attribute.href("https://github.com/dhth/cueitup"),
            attribute.target("_blank"),
          ],
          [element.text("cueitup")],
        ),
      ],
    ),
    consumer_info(config),
    html.button(
      [
        attribute.class(
          "font-semibold px-4 py-1 bg-[#83a598] text-[#282828] hover:bg-[#fabd2f]",
        ),
        attribute.disabled(model.fetching),
        event.on_click(types.FetchMessages(1)),
      ],
      [element.text("Fetch next")],
    ),
    html.button(
      [
        attribute.class(
          "font-semibold px-4 py-1 bg-[#83a598] text-[#282828] hover:bg-[#fabd2f]",
        ),
        attribute.disabled(model.fetching),
        event.on_click(types.FetchMessages(10)),
      ],
      [element.text("Fetch multiple")],
    ),
    html.button(
      [
        attribute.class(
          "font-semibold px-4 py-1 bg-[#bdae93] text-[#282828] hover:bg-[#fabd2f]",
        ),
        attribute.disabled(model.fetching),
        event.on_click(types.ClearMessages),
      ],
      [element.text("Clear Messages")],
    ),
    html.div(
      [
        attribute.class(
          "border-2 border-[#928374] border-opacity-40 border-dashed font-semibold px-4 py-1 flex items-center space-x-4",
        ),
      ],
      [
        html.div([attribute.class("flex items-center space-x-2")], [
          html.label(
            [
              attribute.class("cursor-pointer"),
              attribute.for("hover-control-input"),
            ],
            [element.text("select on hover")],
          ),
          html.input([
            attribute.class(
              "w-4 h-4 text-[#fabd2f] bg-[#282828] focus:ring-[#fabd2f] cursor-pointer",
            ),
            attribute.id("hover-control-input"),
            attribute.type_("checkbox"),
            event.on_check(types.HoverSettingsChanged),
            attribute.checked(model.behaviours.select_on_hover),
          ]),
        ]),
        html.div([attribute.class("flex items-center space-x-2")], [
          html.label(
            [
              attribute.class("cursor-pointer"),
              attribute.for("hover-control-input"),
            ],
            [element.text("delete")],
          ),
          html.input([
            attribute.class(
              "w-4 h-4 text-[#fabd2f] bg-[#282828] focus:ring-[#fabd2f] cursor-pointer",
            ),
            attribute.id("delete-messages"),
            attribute.type_("checkbox"),
            event.on_check(types.DeleteSettingsChanged),
            attribute.checked(model.behaviours.delete_messages),
          ]),
        ]),
        html.div([attribute.class("flex items-center space-x-2")], [
          html.div([attribute.class("relative group")], [
            html.label(
              [
                attribute.class("cursor-pointer"),
                attribute.for("show-live-count"),
              ],
              [element.text("live count")],
            ),
            html.div(
              [
                attribute.class(
                  "absolute left-1/2 -translate-x-1/2 bottom-full mb-2 hidden group-hover:block bg-[#928374] text-[#282828] text-sm px-2 py-1 min-w-[250px]",
                ),
              ],
              [html.text("Message count may fluctuate a bit, that's normal")],
            ),
          ]),
          html.input([
            attribute.class(
              "w-4 h-4 text-[#fabd2f] bg-[#282828] focus:ring-[#fabd2f] cursor-pointer",
            ),
            attribute.id("show-live-count"),
            attribute.type_("checkbox"),
            event.on_check(types.ShowLiveCountChanged),
            attribute.checked(model.behaviours.show_live_count),
          ]),
          case model.message_count, model.behaviours.show_live_count {
            option.Some(c), True ->
              html.p([], [
                element.text(
                  "("
                  <> c
                  |> int.to_string
                  <> " available)",
                ),
              ])
            _, _ -> element.none()
          },
        ]),
      ],
    ),
  ])
}

fn consumer_info(config: Config) -> element.Element(Msg) {
  let profile_name = case config.profile_name |> string.length {
    n if n <= profile_name_max_width -> config.profile_name
    _ -> config.profile_name |> string.slice(0, profile_name_max_width)
  }
  html.div(
    [attribute.class("font-bold px-4 py-1 flex items-center space-x-2")],
    [html.p([attribute.class("text-[#fabd2f]")], [element.text(profile_name)])],
  )
}

fn error_section(model: Model) -> element.Element(Msg) {
  case model.http_error {
    option.None -> element.none()
    option.Some(err) ->
      html.div(
        [
          attribute.role("alert"),
          attribute.class(
            "text-[#fb4934] border-2 border-[#fb4934] border-opacity-50 px-4 py-4 mt-4",
          ),
        ],
        [
          html.strong([attribute.class("font-bold")], [html.text("Error: ")]),
          html.span([attribute.class("block sm:inline")], [
            html.text(err |> http_error_to_string),
          ]),
        ],
      )
  }
}
