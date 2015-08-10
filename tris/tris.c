/* THIS SOFTWARE IS IN THE PUBLIC DOMAIN.
 * WITHOUT ANY WARRANTIES. SO, DON'T FOOL YOUR SELF AND
 * DON'T DO ANY EXPECTATIONS ABOUT IT.
 *
 * AUTHOR: jk <younghead@ukr.net>
 */

/* TODO: PROPER TIMING: FIX THE ULTIMATE PAUSE BY KEY PRESS BUG
 * TODO: CLEAN THE CODE, RENAME THE VARS...
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <termbox.h>

#define NELEMS(x) (sizeof(x)/sizeof((x)[0]))

/* should be called only when termbox is closed (shutdowned) */
void die(const char *fmt, ...)
{
	va_list va;
	va_start(va, fmt);
	vfprintf(stderr, fmt, va);
	va_end(va);
	exit(1);
}

void tb_printf_at(int row, int col, const char *fmt, ...)
{
	int i;
	char buf[1024];
	va_list va;

	va_start(va, fmt);
	vsnprintf(buf, sizeof(buf), fmt, va);
	va_end(va);

	for (i=0; i<sizeof(buf) && buf[i] != '\0'; i++)
		tb_change_cell(col+i, row, buf[i], TB_DEFAULT, TB_DEFAULT);
	tb_present();
}

void tb_printf(const char *fmt, ...)
{
	static int row = 21;
	int i;
	char buf[1024];
	va_list va;

	va_start(va, fmt);
	vsnprintf(buf, sizeof(buf), fmt, va);
	va_end(va);

	for (i=0; i<sizeof(buf) && buf[i] != '\0'; i++)
		tb_change_cell(i, row, buf[i], TB_DEFAULT, TB_DEFAULT);
	row++;
	tb_present();
}

static const unsigned short Tetraminoz[] = {
	0xCC00, 0xCC00, 0xCC00, 0xCC00, 
	0x6C00, 0x8C40, 0x06C0, 0x4620, 
	0xC600, 0x4C80, 0x0C60, 0x2640, 
	0x4444, 0x00F0, 0x2222, 0x0F00,
	0x4E00, 0x4C40, 0x0E40, 0x4640,
	0x4460, 0x2E00, 0xC440, 0x0E80,
	0x44C0, 0x0E20, 0x6440, 0x8E00,
};
static int X = 4;
static int Y = 0;

static unsigned short board[20][10];

unsigned short t, next_t;
unsigned short color, next_color;

unsigned score;
unsigned speed;
unsigned lines_on_this_speed;

unsigned IsColliding(void)
{
	int r, c, b;

	for (b=16; b>0; b--) {
		if (t&(1<<(b-1))) {
			r = (16-b)/4;
			c = (16-b)%4;
			if ((Y+r >= 20) || (Y+r < 0) ||
				(X+c >= 10) || (X+c < 0) ||
				(board[Y+r][X+c])) {

				/*
				 *if (Y == 0 && board[Y+r][X+c])
				 *    GameOver();
				 */

				return 1;
			}
		}
	}

	return 0;
}

void WaitEnter(void)
{
	struct tb_event ev;

	for (;;) {
		tb_poll_event(&ev);
		if (ev.key == TB_KEY_ENTER) 
			break;
	}
}

void GameOver(void)
{
	tb_clear();

	tb_printf_at(1, 2, "GAME OVER!");

	tb_printf_at(3, 2, "YOUR SCORE IS %d", score);
	tb_printf_at(4, 2, "YOU FINISHED AT SPEED %d", speed);

	tb_printf_at(9, 2, "<hit Enter>");
	tb_present();

	WaitEnter();
	tb_shutdown();
	exit(0);
}

void YouWin(void)
{
	tb_clear();

	tb_printf_at(1, 2, "CONGRATULATIONS! YOU WIN!");

	tb_printf_at(3, 2, "YOU HAVE COMPLETED THE GAME WITH SCORE %d", score);

	tb_printf_at(9, 2, "<hit Enter>");
	tb_present();

	WaitEnter();
	tb_shutdown();
	exit(0);
}

void Pause(void)
{
	tb_clear();

	tb_printf_at(1, 2, "... PAUSE ...");

	tb_printf_at(9, 2, "<hit Enter>");
	tb_present();

	WaitEnter();
}

void New(void)
{
	static unsigned short colors[] = {
		TB_RED, TB_GREEN, TB_BLUE, TB_MAGENTA, TB_CYAN
	};

	if (next_t == 0) {
		do {
			next_t = Tetraminoz[rand() % NELEMS(Tetraminoz)];
		} while ((next_t & 0xF000) == 0);

		next_color = colors[rand() % NELEMS(colors)];
		if (rand() % 2)
			next_color |= TB_BOLD;
	}

	t = next_t;
	do {
		next_t = Tetraminoz[rand() % NELEMS(Tetraminoz)];
	} while ((next_t & 0xF000) == 0);

	color = next_color;
	next_color = colors[rand() % NELEMS(colors)];
	if (rand() % 2)
		next_color |= TB_BOLD;

	X = 4;
	Y = 0;

	if (IsColliding())
		GameOver();
}

void Draw(void)
{
	int r, c;
	tb_clear();
	for (r = 0; r < 20; r++)
		for (c = 0; c < 10; c++) {
			if (board[r][c]) {
				tb_change_cell(2*c, r, ' ', TB_DEFAULT, board[r][c]);
				tb_change_cell(2*c+1, r, ' ', TB_DEFAULT, board[r][c]);
			} else {
				tb_change_cell(2*c, r, ' ', TB_DEFAULT, TB_WHITE);
				tb_change_cell(2*c+1, r, ' ', TB_DEFAULT, TB_WHITE);
			}
		}

	for (r = 0; r < 4; r++)
		for (c = 0; c < 4; c++)
			if ((1<<(15-(r*4+c))) & t) {
				tb_change_cell(2*(X+c), Y+r, ' ', TB_DEFAULT, color);
				tb_change_cell(2*(X+c)+1, Y+r, ' ', TB_DEFAULT, color);
			}

	/* next tetramino */
	for (r = 0; r < 4; r++)
		for (c = 0; c < 4; c++)
			if ((1<<(15-(r*4+c))) & next_t) {
				tb_change_cell(2*(15+c), 2+r, ' ', TB_DEFAULT, next_color);
				tb_change_cell(2*(15+c)+1, 2+r, ' ', TB_DEFAULT, next_color);
			}

	tb_printf_at(22, 0, "score: %d\n", score);
	tb_printf_at(23, 0, "speed: %d\n", speed);
	tb_present();
}

void Rotate(void)
{
	unsigned short oldt = t;
	int i;

	for (i = 0; i < NELEMS(Tetraminoz); i++)
		if (t == Tetraminoz[i])
			break;
	if ((i & 3) == 3)
		t =  Tetraminoz[i & ~3];
	else
		t =  Tetraminoz[i+1];

	if (IsColliding())
		t = oldt;
}

void Release(void)
{
	int r, c, b;

	for (b=16; b>0; b--) {
		if (t&(1<<(b-1))) {
			r = (16-b)/4;
			c = (16-b)%4;
			board[Y+r][X+c] = color;
			}
		}

	t = 0;
}

void EliminateLines(void)
{
	int consecutive_lines = 0;
	int r, rr, c;

	for (r=19; r>=0; r--) {
		for (c=0; c<10; c++) {
			if (board[r][c] == 0)
				break;
		}

		if (c == 10) {
			consecutive_lines++;
			for (rr=r; rr>0; rr--)
				for (c=0; c<10; c++)
					board[rr][c] = board[rr-1][c];
			r++; /* need to stay on the same line */
		}
	}

	/* FIXME: maybe, it's wrong to calculate 3 lines a 7 if it is
	 * actually 4, but has some hole on some line. Maybe, it should be
	 * calculated as 1 + 3 ? */
	switch (consecutive_lines) {
		case 1: score += 1; break;
		case 2: score += 3; break;
		case 3: score += 7; break;
		case 4: score += 15; break;
	}

	lines_on_this_speed += consecutive_lines;
	if (lines_on_this_speed > 10) {
		speed++;
		lines_on_this_speed = 0;

		if (speed == 10)
			YouWin();
	}
}

int Fall(void)
{
	Y++;
	if (IsColliding()) {
		Y--;
		Release();
		EliminateLines();
		New();
		return 0;
	}
	return 1;
}

void Fallthrough(void)
{
	while (Fall())
		;
}

int Left(void)
{
	X--;
	if (IsColliding()) {
		X++;
		return 0;
	}
	return 1;
}

int Right(void)
{
	X++;
	if (IsColliding()) {
		X--;
		return 0;
	}
	return 1;
}

int main(void)
{
	int ret;
	struct tb_event ev;

	if (tb_init() < 0)
		die("Failed to initialize termbox. Exiting...\n");

	srand(time(NULL));

	New();

	for (;;) {
		Draw();

		/* FIXME: check what is returned */
		ret = tb_peek_event(&ev, 1000-100*speed);
		if (ret == 0) {
			Fall();
		} else if (ev.type == TB_EVENT_KEY) {
			if (ev.key == TB_KEY_ARROW_UP || ev.ch == 'k') {
				Rotate();
			} else if (ev.ch == 'K') {
				Rotate();
				Rotate();
			} else if (ev.key == TB_KEY_ARROW_DOWN || ev.ch == 'j') {
				Fall();
			} else if (ev.key == TB_KEY_ARROW_LEFT || ev.ch == 'h') {
				Left();
			} else if (ev.ch == 'H') {
				while (Left())
					;
			} else if (ev.key == TB_KEY_ARROW_RIGHT || ev.ch == 'l') {
				Right();
			} else if (ev.ch == 'L') {
				while (Right())
					;
			} else if (ev.ch == 'p' || ev.ch == 'P') {
				Pause();
			} else if (ev.key == TB_KEY_ESC || ev.ch == 'q') {
				break;
			} else if (ev.key == TB_KEY_SPACE || ev.ch == ' ') {
				Fallthrough();
			}
		}
	}

	tb_shutdown();
}
