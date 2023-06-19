// TODO:
// generate new random food location
// check that food location is not occupied by the snake

use sdl2::pixels::Color;
use crate::{coord::Coord, config::Config, renderer::Renderer};

use rand::Rng;

pub struct Food(pub Coord);

impl Food {
    // TODO: check that there is no cells of snake
    pub fn new_random(config: &Config) -> Self {
        let mut rng = rand::thread_rng();

        Self(Coord(
            (rng.gen::<u32>() % config.width) as i32,
            (rng.gen::<u32>() % config.height) as i32,
        ))
    }

    pub fn draw(&self, renderer: &mut Renderer) {
        renderer.set_color(Color::RED);
        renderer.draw_dot(self.0);
    }
}
