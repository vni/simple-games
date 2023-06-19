use crate::{config::Config, coord::Coord};
use sdl2::{pixels::Color, rect::Rect, render::WindowCanvas, video::Window};

pub struct Renderer<'a> {
    canvas: WindowCanvas,
    config: &'a Config,
}

impl<'a> Renderer<'a> {
    pub fn new(window: Window, config: &'a Config) -> Self {
        Self {
            canvas: window
                .into_canvas()
                .build()
                .expect("Failed to transform window into canvas"),
            config,
        }
    }

    pub fn draw_background(&mut self) {
        self.canvas.set_draw_color(Color::RGB(0x20, 0x30, 0x40));
        self.canvas.clear();
    }

    pub fn draw_dot(&mut self, coord: Coord) {
        self.canvas
            .fill_rect(Rect::new(
                coord.0 * self.config.pixel_size as i32,
                coord.1 * self.config.pixel_size as i32,
                self.config.pixel_size,
                self.config.pixel_size,
            ))
            .expect("Failed to fill_rect");
    }

    pub fn set_color(&mut self, color: Color) {
        self.canvas.set_draw_color(color);
    }

    pub fn present(&mut self) {
        self.canvas.present();
    }

    pub fn get_width(&self) -> u32 {
        self.config.width
    }

    pub fn get_height(&self) -> u32 {
        self.config.height
    }
}
