import type { IndexProps } from "~/types.page.about.gen.d.ts";
import Header from "~/components/Header";
import "../about.css";


const  IndexPage = (p: IndexProps) => {
  return (
    <>
      <Header />
      <div className="about-page">
        <div className="about-breadcrumb">
          <a href="/about">About</a>
          <span className="breadcrumb-sep">/</span>
          <span>Nested</span>
        </div>

        <h1 className="about-title">About / Nested</h1>
        <p className="about-message">{p.message}</p>

        <div className="about-card">
          <h2>親ルート</h2>
          <p>このページは <code>/about/nested</code> にあります。</p>
          <p>
            <a href="/about/nested/nested" className="about-link">
              /about/nested/nested を見る →
            </a>
          </p>
          <p>
            <a href="/about" className="about-link">
              ← /about に戻る
            </a>
          </p>
        </div>
      </div>
    </>
  );
};


export default IndexPage;
