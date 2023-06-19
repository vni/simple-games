use sdl2::{event::Event, keyboard::Keycode};

mod config;
mod coord;
mod food;
//mod game_world;
mod renderer;
mod snake;

fn main() {
    let config = config::Config::from_args();
    println!("config: {}", config);

    let sdl_context = sdl2::init().expect("Failed to init sdl2");
    let video_subsystem = sdl_context.video().expect("Failed to get video_subsystem");
    let window = video_subsystem
        .window(
            "snake game",
            config.width * config.pixel_size,
            config.height * config.pixel_size,
        )
        .opengl()
        .build()
        .expect("Failed to create application window");

    //let mut game_world = game_world::GameWorld::new();
    let mut snake = snake::Snake::new(config.width, config.height);
    let mut food = food::Food::new_random(&config);
    let mut renderer = renderer::Renderer::new(window, &config);

    let mut event_pump = sdl_context
        .event_pump()
        .expect("Failed to get event_pump from sdl_context");
    let mut speed = 100;

    let mut frame_counter = 0;
    'game_loop: loop {
        for event in event_pump.poll_iter() {
            match event {
                Event::Quit { .. } => break 'game_loop,
                Event::KeyDown {
                    keycode: Some(Keycode::Escape),
                    ..
                } => break 'game_loop,
                Event::KeyDown {
                    keycode: Some(Keycode::Q),
                    ..
                } => break 'game_loop,
                Event::KeyDown {
                    keycode: Some(Keycode::Up),
                    ..
                } => snake.turn_up(),
                Event::KeyDown {
                    keycode: Some(Keycode::Down),
                    ..
                } => snake.turn_down(),
                Event::KeyDown {
                    keycode: Some(Keycode::Left),
                    ..
                } => snake.turn_left(),
                Event::KeyDown {
                    keycode: Some(Keycode::Right),
                    ..
                } => snake.turn_right(),
                Event::KeyDown {
                    keycode: Some(Keycode::Num9),
                    ..
                } => { speed -= 5; println!("speed: {speed}"); }
                Event::KeyDown {
                    keycode: Some(Keycode::Num0),
                    ..
                } => { speed += 5; println!("speed: {speed}"); }
                _ => {/* println!("Unknown key: {event:?}"); */}
            }
        }

        //game_world.draw(&mut renderer);
        renderer.draw_background();
        snake.draw(&mut renderer);
        food.draw(&mut renderer);
        renderer.present();

        std::thread::sleep(std::time::Duration::new(0, 1_000_000_000u32 / (30 + snake.apples_eaten() + speed - 100)));
        frame_counter += 1;
        if frame_counter % 10 == 0 {
            snake.tick();
            if *snake.head() == food.0 {
                snake.eat_apple();
                food = food::Food::new_random(&config);
            }
            frame_counter = 0;
        }
    }
}
