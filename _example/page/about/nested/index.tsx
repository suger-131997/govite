import type { IndexProps } from "~/types.gen.page.about.nested.d.ts";
import Header from "~/components/Header";
import "../../about.css";

const IndexPage = (p: IndexProps) => {
  return (
    <>
      <Header />
      <div className="about-page">
        <div className="about-breadcrumb">
          <a href="/about">About</a>
          <span className="breadcrumb-sep">/</span>
          <a href="/about/nested">Nested</a>
          <span className="breadcrumb-sep">/</span>
          <span>Nested</span>
        </div>

        <h1 className="about-title">About / Nested / Nested</h1>
        <p className="about-message">{p.message}</p>

        <div className="about-card">
          <h2>親ルート</h2>
          <p>このページは <code>/about/nested/nested</code> にあります。</p>
          <p>
            <a href="/about/nested" className="about-link">
              ← /about/nested に戻る
            </a>
          </p>
        </div>
      </div>
    </>
  );
};

export default IndexPage;
