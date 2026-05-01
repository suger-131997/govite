import type { NotFoundProps } from "~/types.gen.ts";
import Header from "~/components/Header";
import "./not_found.css";

const NotFoundPage = (p: NotFoundProps) => {
  return (
    <>
      <Header />
      <div className="not-found-page">
        <div className="not-found-container">
          <p className="not-found-status">404</p>
          <h1 className="not-found-title">ページが見つかりません</h1>
          <p className="not-found-description">
            <code className="not-found-path">{p.path}</code>{" "}
            は存在しないか、移動された可能性があります。
          </p>
          <div className="not-found-actions">
            <a href="/" className="btn-primary">ホームへ戻る</a>
            <a href="/users" className="btn-secondary">Users を見る</a>
          </div>
        </div>
      </div>
    </>
  );
};

export default NotFoundPage;
