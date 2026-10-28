#[path = "../rust/kartouche.rs"] mod kartouche;
use kartouche::{Outcome, Style};

fn main() {
    eprintln!("Loading configuration...");
    kartouche::banner(None);
    eprintln!("Checking inputs...");
    kartouche::card();
    eprintln!("Reviewing results...");
    kartouche::summary(Style::Rail, &[("connect", Outcome::Ok), ("fetch", Outcome::Failed), ("deploy", Outcome::Skipped)]);
    eprintln!("Continuing...");
}
