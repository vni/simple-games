use crate::{renderer::Renderer, snake::Snake};

pub struct GameWorld {
    snake: Snake,
}

impl GameWorld {
    pub fn new() -> GameWorld {
        GameWorld {
            snake: Snake::new(),
        }
    }

    pub fn draw(&mut self, renderer: &mut Renderer) {
        // background
        renderer.draw_background();

        // snake
        self.snake.draw(renderer);

        // food

        renderer.present();
    }
}
