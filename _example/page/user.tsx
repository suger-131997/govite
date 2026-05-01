import type { UserDetailProps } from "~/types.gen.page.d.ts";
import Header from "~/components/Header";
import "./user.css";

const roleColors: Record<string, string> = {
  Admin: "#7c3aed",
  Editor: "#0891b2",
  Viewer: "#059669",
  Developer: "#d97706",
  Manager: "#dc2626",
};

const statusClass: Record<string, string> = {
  Active: "status-active",
  Inactive: "status-inactive",
  Pending: "status-pending",
};

const UserPage = (p: UserDetailProps) => {
  if (!p.found) {
    return (
      <>
        <Header />
        <div className="user-page">
          <div className="not-found">
            <p className="not-found-code">404</p>
            <p className="not-found-msg">User not found</p>
            <a href="/users" className="back-link">← Users に戻る</a>
          </div>
        </div>
      </>
    );
  }

  const { user } = p;

  return (
    <>
      <Header />
      <div className="user-page">
        <div className="user-breadcrumb">
          <a href="/users">Users</a>
          <span className="breadcrumb-sep">/</span>
          <span>{user.name}</span>
        </div>

        <div className="user-card">
          <div className="user-avatar">
            {user.name.split(" ").map((w) => w[0]).join("").slice(0, 2)}
          </div>

          <div className="user-card-body">
            <h1 className="user-name">{user.name}</h1>
            <div className="user-meta">
              <span
                className="role-badge"
                style={{ backgroundColor: roleColors[user.role] ?? "#6b7280" }}
              >
                {user.role}
              </span>
              <span className={`status-badge ${statusClass[user.status] ?? ""}`}>
                {user.status}
              </span>
            </div>
          </div>
        </div>

        <div className="detail-grid">
          <div className="detail-section">
            <h2 className="section-title">基本情報</h2>
            <dl className="detail-list">
              <div className="detail-row">
                <dt>ID</dt>
                <dd>#{user.id}</dd>
              </div>
              <div className="detail-row">
                <dt>名前</dt>
                <dd>{user.name}</dd>
              </div>
              <div className="detail-row">
                <dt>メールアドレス</dt>
                <dd>
                  <a href={`mailto:${user.email}`} className="email-link">
                    {user.email}
                  </a>
                </dd>
              </div>
              <div className="detail-row">
                <dt>ロール</dt>
                <dd>{user.role}</dd>
              </div>
              <div className="detail-row">
                <dt>ステータス</dt>
                <dd>
                  <span className={`status-badge ${statusClass[user.status] ?? ""}`}>
                    {user.status}
                  </span>
                </dd>
              </div>
              <div className="detail-row">
                <dt>登録日</dt>
                <dd>{user.joinedAt}</dd>
              </div>
            </dl>
          </div>
        </div>

        <div className="user-actions">
          <a href="/users" className="back-link">← Users に戻る</a>
        </div>
      </div>
    </>
  );
};

export default UserPage;
