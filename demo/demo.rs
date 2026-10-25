#[path = "../rust/kartouche.rs"] mod kartouche;
use kartouche::{Style, Stream};

fn main() {
    use kartouche as k;
    let out = Stream::Stdout;
    println!("{}", k::render_banner("app_name", "0.1.0", "app_description", None, out));
    println!();
    println!("{}", k::render_card("lumen", "0.3.1", "tiny notes, zero friction", out));
    for (name, style) in [("inline", Style::Inline), ("dots", Style::Dots), ("crumbs", Style::Crumbs),
                          ("rail", Style::Rail), ("line", Style::Line), ("log", Style::Log)] {
        println!("\n-- {name}");
        println!("{}", k::render_title("app_name", "0.1.0", out));
        println!("{}", k::render_frame(style, 3, 4, "build", &["connect", "fetch", "build", "deploy"], out));
    }
    println!("\n-- end");
    println!("{}", k::render_end(true, "done in 4.2s", out));
    println!("{}", k::render_end(false, "failed at 3/4 build", out));
}
