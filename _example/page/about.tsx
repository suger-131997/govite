import type { AboutProps } from "~/types.gen.page.d.ts";
import Header from "~/components/Header";
import "./about.css";

const AboutPage = (p: AboutProps) => {
  return (
    <>
      <Header />
      <div className="about-page">
        <div className="about-breadcrumb">
          <span>About</span>
        </div>

        <h1 className="about-title">About</h1>
        <p className="about-message">{p.message}</p>

        <div className="about-card">
          <h2>ネストされたルート</h2>
          <p>このページは <code>/about</code> にあります。</p>
          <p>
            <a href="/about/IndexPage" className="about-link">
              /about/nested を見る →
            </a>
          </p>
        </div>
      </div>
    </>
  );
};

export default AboutPage;
