import "./Header.css";

const navItems = [
  { href: "/", label: "Home" },
  { href: "/users", label: "Users" },
];

const Header = () => {
  const current = typeof window !== "undefined" ? window.location.pathname : "";

  return (
    <header className="site-header">
      <div className="site-header-inner">
        <span className="site-logo">govite</span>
        <nav className="site-nav">
          {navItems.map(({ href, label }) => (
            <a
              key={href}
              href={href}
              className={`nav-link${current === href ? " nav-link-active" : ""}`}
            >
              {label}
            </a>
          ))}
        </nav>
      </div>
    </header>
  );
};

export default Header;
