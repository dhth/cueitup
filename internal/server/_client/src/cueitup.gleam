import effects.{fetch_behaviours, fetch_config}
import lustre
import lustre/effect
import model.{type Model, init_model}
import types.{type Msg}
import update
import view

pub fn main() {
  let app = lustre.application(init, update.update, view.view)
  let assert Ok(_) = lustre.start(app, "#app", Nil)
}

fn init(_) -> #(Model, effect.Effect(Msg)) {
  #(init_model(), effect.batch([fetch_config(), fetch_behaviours()]))
}
