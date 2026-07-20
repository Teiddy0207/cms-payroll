service phải gọi đến mapper
controller khai báo query params và dùng ở repository, k dùng ở servivce

module nào thì phải chia ra, ví dụ job position chỉ thuộc file job posiion thôi, k để chung

5. Các tham số lọc QueryParams (ví dụ lọc theo phòng ban department_id) phải được định nghĩa ở Controller và truyền thẳng xuống Repository để dựng SQL động, tuyệt đối không bóc tách hoặc xử lý logic QueryParams ở tầng Service.