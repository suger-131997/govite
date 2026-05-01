import type { IndexProps } from "~/types.gen.ts";
import bannerImg from "~/assets/banner.png";
import goLogo from "~/assets/go-logo.svg";
import viteLogo from "~/assets/vite-logo.svg";
import "./index.css";

const IndexPage = (p: IndexProps) => {
  return (
    <div className="index-page">
      <img src={bannerImg} className="banner" alt="govite banner" />

      <div className="header">
        <div className="logos">
          <img src={goLogo} className="logo" alt="Go" />
          <span className="logo-plus">+</span>
          <img src={viteLogo} className="logo" alt="Vite" />
        </div>
        <h1>Complex Props Demo</h1>
        <p className="subtitle">
          Demonstrating type-safe props passing from Go to React
        </p>
      </div>

      {/* Profile: string, int, float64, bool, *string */}
      <div className="card">
        <h2 className="card-title">Profile</h2>
        <div className="profile-grid">
          <div className="field">
            <span className="field-label">Name</span>
            <span className="field-value">{p.name}</span>
          </div>
          <div className="field">
            <span className="field-label">Age</span>
            <span className="field-value">{p.age}</span>
          </div>
          <div className="field">
            <span className="field-label">Score</span>
            <span className="field-value">{p.score}</span>
          </div>
          <div className="field">
            <span className="field-label">Status</span>
            <span
              className={`field-value ${p.isActive ? "status-active" : "status-inactive"}`}
            >
              {p.isActive ? "Active" : "Inactive"}
            </span>
          </div>
          <div className="field">
            <span className="field-label">Website</span>
            <span className="field-value">
              {p.website ? (
                <a href={p.website} target="_blank" rel="noreferrer">
                  {p.website}
                </a>
              ) : (
                "—"
              )}
            </span>
          </div>
        </div>
      </div>

      {/* Address: nested struct */}
      <div className="card">
        <h2 className="card-title">Address</h2>
        <div className="address-text">
          {p.address.street}
          <br />
          {p.address.city} {p.address.zipCode}
        </div>
      </div>

      {/* Tags: []struct */}
      <div className="card">
        <h2 className="card-title">Tags</h2>
        <div className="tags-list">
          {p.tags.map((tag) => (
            <span key={tag.name} className="badge">
              <span
                className="badge-dot"
                style={{ backgroundColor: tag.color }}
              />
              {tag.name}
            </span>
          ))}
        </div>
      </div>

      {/* Metadata: map[string]string */}
      <div className="card">
        <h2 className="card-title">Metadata</h2>
        <table className="meta-table">
          <tbody>
            {Object.entries(p.metadata).map(([key, value]) => (
              <tr key={key}>
                <td>{key}</td>
                <td>{value}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* CreatedAt: time.Time */}
      <div className="card">
        <h2 className="card-title">Timestamp</h2>
        <div className="field">
          <span className="field-label">Created At</span>
          <span className="timestamp">{p.createdAt}</span>
        </div>
      </div>
    </div>
  );
};

export default IndexPage;
