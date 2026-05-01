import type { UsersProps } from "~/types.gen.ts";
import "./users.css";

const statusClass: Record<string, string> = {
  Active: "status-active",
  Inactive: "status-inactive",
  Pending: "status-pending",
};

function buildPageUrl(page: number, pageSize: number): string {
  return `?page=${page}&pageSize=${pageSize}`;
}

const UsersPage = (p: UsersProps) => {
  const { users, currentPage, totalPages, totalUsers, pageSize } = p;

  const pages: number[] = [];
  const delta = 2;
  for (let i = 1; i <= totalPages; i++) {
    if (
      i === 1 ||
      i === totalPages ||
      (i >= currentPage - delta && i <= currentPage + delta)
    ) {
      pages.push(i);
    }
  }

  const pageButtons: (number | "ellipsis")[] = [];
  for (let i = 0; i < pages.length; i++) {
    if (i > 0 && pages[i] - pages[i - 1] > 1) {
      pageButtons.push("ellipsis");
    }
    pageButtons.push(pages[i]);
  }

  const startRow = (currentPage - 1) * pageSize + 1;
  const endRow = Math.min(currentPage * pageSize, totalUsers);

  return (
    <div className="users-page">
      <header className="users-header">
        <h1>Users</h1>
        <p className="users-subtitle">
          Showing {startRow}–{endRow} of {totalUsers} users
        </p>
      </header>

      <div className="table-wrapper">
        <table className="users-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Status</th>
              <th>Joined At</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td className="col-id">{user.id}</td>
                <td className="col-name">{user.name}</td>
                <td className="col-email">{user.email}</td>
                <td>{user.role}</td>
                <td>
                  <span className={`status-badge ${statusClass[user.status] ?? ""}`}>
                    {user.status}
                  </span>
                </td>
                <td className="col-date">{user.joinedAt}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="pagination">
        <a
          href={buildPageUrl(1, pageSize)}
          className={`page-btn${currentPage === 1 ? " disabled" : ""}`}
          aria-disabled={currentPage === 1}
        >
          «
        </a>
        <a
          href={buildPageUrl(Math.max(1, currentPage - 1), pageSize)}
          className={`page-btn${currentPage === 1 ? " disabled" : ""}`}
          aria-disabled={currentPage === 1}
        >
          ‹
        </a>

        {pageButtons.map((btn, i) =>
          btn === "ellipsis" ? (
            <span key={`ellipsis-${i}`} className="page-ellipsis">
              …
            </span>
          ) : (
            <a
              key={btn}
              href={buildPageUrl(btn, pageSize)}
              className={`page-btn${btn === currentPage ? " active" : ""}`}
              aria-current={btn === currentPage ? "page" : undefined}
            >
              {btn}
            </a>
          )
        )}

        <a
          href={buildPageUrl(Math.min(totalPages, currentPage + 1), pageSize)}
          className={`page-btn${currentPage === totalPages ? " disabled" : ""}`}
          aria-disabled={currentPage === totalPages}
        >
          ›
        </a>
        <a
          href={buildPageUrl(totalPages, pageSize)}
          className={`page-btn${currentPage === totalPages ? " disabled" : ""}`}
          aria-disabled={currentPage === totalPages}
        >
          »
        </a>

        <span className="page-size-label">表示件数:</span>
        {[5, 10, 20].map((size) => (
          <a
            key={size}
            href={buildPageUrl(1, size)}
            className={`page-btn${size === pageSize ? " active" : ""}`}
          >
            {size}
          </a>
        ))}
      </div>
    </div>
  );
};

export default UsersPage;
