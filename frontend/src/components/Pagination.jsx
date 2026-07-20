export function Pagination({ page, totalItems, pageSize, onPageChange }) {
  const totalPages = Math.ceil(totalItems / pageSize) || 1;
  const start = (page - 1) * pageSize + 1;
  const end = Math.min(page * pageSize, totalItems);

  const getPages = () => {
    const pages = [];
    const maxVisible = 5;
    let startPage = Math.max(1, page - Math.floor(maxVisible / 2));
    let endPage = Math.min(totalPages, startPage + maxVisible - 1);
    if (endPage - startPage < maxVisible - 1) {
      startPage = Math.max(1, endPage - maxVisible + 1);
    }
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }
    return pages;
  };

  if (totalItems === 0) return null;

  return (
    <div className="pagination">
      <div className="pagination-info">
        Hiển thị <strong>{start}–{end}</strong> trong <strong>{totalItems}</strong> kết quả
      </div>
      <div className="pagination-controls">
        <button
          className="page-btn"
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          aria-label="Trang trước"
        >
          ‹
        </button>
        {getPages().map((p) => (
          <button
            key={p}
            className={`page-btn${page === p ? ' active' : ''}`}
            onClick={() => onPageChange(p)}
          >
            {p}
          </button>
        ))}
        <button
          className="page-btn"
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          aria-label="Trang sau"
        >
          ›
        </button>
      </div>
    </div>
  );
}

export default Pagination;
