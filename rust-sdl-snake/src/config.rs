pub struct Config {
    pub width: u32,
    pub height: u32,
    pub pixel_size: u32,
}

impl Config {
    pub fn from_args() -> Self {
        /*
        let args = std::env::args().collect();

        for arg in args.iter() {
        }
        */

        // TODO: Default values. Let them be readed from command line.
        Config {
            width: 80,
            height: 40,
            pixel_size: 20,
        }
    }
}

impl std::fmt::Display for Config {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(
            f,
            "Config {{ width: {}, height: {}, pixel_size: {} }}",
            self.width, self.height, self.pixel_size
        )
    }
}
