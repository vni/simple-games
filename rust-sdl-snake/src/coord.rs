#[derive(Clone, Copy, PartialEq)]
pub struct Coord(pub i32, pub i32);

impl std::ops::Add for Coord {
    type Output = Self;

    fn add(self, other: Self) -> Self::Output {
        Self(self.0 + other.0, self.1 + other.1)
    }
}
