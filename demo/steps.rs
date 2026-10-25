#[path = "../rust/kartouche.rs"] mod kartouche;
use kartouche::{Steps, Style};
use std::{env, thread, time::Duration};

fn main() {
    let style = match env::args().nth(1).as_deref() {
        Some("dots") => Style::Dots, Some("crumbs") => Style::Crumbs,
        Some("rail") => Style::Rail, Some("line") => Style::Line,
        Some("log") => Style::Log, _ => Style::Inline,
    };
    let mut s = Steps::start(style, &["connect", "fetch", "build", "deploy"]);
    for _ in 0..3 { thread::sleep(Duration::from_millis(300)); s.next(); }
    thread::sleep(Duration::from_millis(300));
    s.done(None);
}
