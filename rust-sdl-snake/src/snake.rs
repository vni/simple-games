// TODO: add clockwise and couterclockwise movement

use crate::{coord::Coord, renderer::Renderer};
use sdl2::pixels::Color;

#[derive(PartialEq)]
enum Direction {
    Up,
    Down,
    Left,
    Right,
}

pub struct Snake {
    body: Vec<Coord>,
    len: usize,
    dir: Direction,
    apples_eaten: u32,
    world_width: u32,
    world_height: u32,
}

impl Snake {
    pub fn new(world_width: u32, world_height: u32) -> Snake {
        Snake {
            body: vec![Coord(10, 10), Coord(11, 10), Coord(12, 10)],
            len: 3,
            dir: Direction::Right,
            apples_eaten: 0,
            world_width,
            world_height,
        }
    }

    pub fn draw(&self, renderer: &mut Renderer) {
        let mut idx = self.body.len();

        renderer.set_color(Color::BLUE);
        idx -= 1;
        renderer.draw_dot(self.body[idx]);

        renderer.set_color(Color::GREEN);
        for _ in 1..self.len {
            idx -= 1;
            renderer.draw_dot(self.body[idx]);
        }
    }

    pub fn head(&self) -> &Coord {
        let idx = self.body.len() - 1;
        &self.body[idx]
    }

    pub fn tick(&mut self) {
        let head = self.body.last().unwrap();
        let mut new_head = match self.dir {
            Direction::Up => *head + Coord(0, -1),
            Direction::Down => *head + Coord(0, 1),
            Direction::Left => *head + Coord(-1, 0),
            Direction::Right => *head + Coord(1, 0),
        };

        if new_head.0 >= self.world_width as i32 {
            new_head.0 = 0;
        }

        if new_head.0 < 0 {
            new_head.0 = (self.world_width - 1) as i32;
        }

        if new_head.1 >= self.world_height as i32 {
            new_head.1 = 0;
        }

        if new_head.1 < 0 {
            new_head.1 = (self.world_height - 1) as i32;
        }


        self.body.push(new_head);

        // FIXME: FINISH HERE. Shrink array.
        /*
        if (self.body.len - self.len) > 1000 {
            self.body ...
        }
        */
    }

    pub fn turn_up(&mut self) {
        /*if self.dir == Direction::Down {
            panic!("OOPS, you crashed on yourself");
        }*/

        self.dir = Direction::Up;
    }

    pub fn turn_down(&mut self) {
        /*if self.dir == Direction::Up {
            panic!("OOPS, you crashed on yourself");
        }*/

        self.dir = Direction::Down;
    }

    pub fn turn_right(&mut self) {
        /*if self.dir == Direction::Left {
            panic!("OOPS, you crashed on yourself");
        }*/

        self.dir = Direction::Right;
    }

    pub fn turn_left(&mut self) {
        /*if self.dir == Direction::Right {
            panic!("OOPS, you crashed on yourself");
        }*/

        self.dir = Direction::Left;
    }

    pub fn eat_apple(&mut self) {
        self.len += 1;
        self.apples_eaten += 1;
        println!("apples_eaten: {}", self.apples_eaten);
    }

    pub fn apples_eaten(&self) -> u32 {
        self.apples_eaten
    }
}
