Discord: Real-Time Architecture at Internet Scale
Discord trông có vẻ chỉ là một ứng dụng chat đơn giản. Bạn gửi một tin nhắn, vài mili giây sau bạn bè của bạn đã nhận được. Nhưng phía sau trải nghiệm gần như tức thời đó là một trong những hệ thống realtime phức tạp nhất Internet, phải xử lý hàng tỷ sự kiện mỗi ngày bao gồm tin nhắn, voice, video, trạng thái online và vô số tương tác khác.
Điều đáng chú ý là Discord không cố gắng xây dựng mọi thứ thật phức tạp ngay từ đầu. Họ bắt đầu với một kiến trúc khá đơn giản, sau đó liên tục thay đổi khi quy mô người dùng tăng lên. Đây cũng là một trong những lý do Discord thường được nhắc đến như một ví dụ điển hình về việc scale từng bước thay vì over-engineering ngay từ ngày đầu.
Bắt đầu với một stack rất đơn giản
Discord ra mắt vào năm 2015 với một stack khá "thực dụng":
Frontend sử dụng React
Backend viết bằng Elixir chạy trên Erlang VM (BEAM)
PostgreSQL lưu trữ dữ liệu
WebSocket cho toàn bộ giao tiếp realtime
Bare Metal thay vì Cloud
Nhìn qua thì đây không phải là stack "thời thượng" nhất thời điểm đó. Tuy nhiên, từng lựa chọn đều phục vụ trực tiếp cho bài toán realtime mà Discord muốn giải quyết.
Vì sao Discord chọn Elixir?
Đây có lẽ là quyết định kỹ thuật nổi tiếng nhất của Discord.
Elixir chạy trên BEAM VM của Erlang – nền tảng vốn được sinh ra để xây dựng các hệ thống viễn thông với yêu cầu uptime cực cao.
BEAM có một số đặc điểm rất phù hợp với ứng dụng chat:
Hỗ trợ hàng triệu lightweight process
Scheduler tối ưu cho concurrency
Isolation giữa các process
Fault tolerance rất cao
Hot code reloading
Thay vì dùng thread truyền thống, Discord gần như tạo một process cho mỗi client kết nối. Điều này nghe có vẻ rất tốn tài nguyên nếu nhìn theo tư duy của Java hay Go, nhưng trên BEAM, mỗi process chỉ chiếm rất ít bộ nhớ nên việc quản lý hàng triệu kết nối đồng thời vẫn hoàn toàn khả thi.
Điểm thú vị là BEAM áp dụng triết lý "Let it crash".
Nếu một process gặp lỗi, hệ thống không cố gắng xử lý mọi exception. Process đó đơn giản bị kill và supervisor sẽ tạo một process mới. Điều này giúp hệ thống ổn định hơn rất nhiều khi vận hành ở quy mô lớn.
WebSocket Gateway – Trái tim của Discord
Khác với HTTP request/response truyền thống, Discord cần duy trì kết nối lâu dài với từng người dùng.
Giải pháp là WebSocket.
Mỗi client mở một kết nối persistent tới Gateway. Tất cả các sự kiện như:
gửi tin nhắn
typing indicator
emoji reaction
online/offline status
voice event
đều đi qua tầng Gateway này. Nhưng hàng triệu WebSocket connection không thể chạy trên một server duy nhất. Discord cần một cách để chia nhỏ toàn bộ hệ thống.
Sharding – Bí quyết giúp Discord scale ngang
Đây là kỹ thuật quan trọng nhất trong toàn bộ kiến trúc. Thay vì có một Gateway khổng lồ, Discord chia hệ thống thành hàng nghìn shard. Mỗi shard chỉ quản lý khoảng vài nghìn kết nối (khoảng 5.000 user theo bài viết). Việc shard mang lại nhiều lợi ích cùng lúc.
Nếu một shard bị lỗi, chỉ một nhóm nhỏ người dùng bị ảnh hưởng thay vì toàn bộ hệ thống.
Memory của từng node cũng dễ kiểm soát hơn vì mỗi server chỉ phải giữ trạng thái của một lượng kết nối nhất định. Khi cần mở rộng, Discord chỉ cần tạo thêm shard mới rồi phân phối người dùng sang các shard đó.
Điều quan trọng là mỗi user hoặc guild luôn được ánh xạ tới một shard xác định thông qua ID hashing. Điều này giúp việc định tuyến request nhanh và ổn định hơn.
Database cũng phải Shard
Ban đầu, Discord lưu toàn bộ dữ liệu trong PostgreSQL. Bao gồm:
User
Guild (Server)
Message
Role
Permission
Khi số lượng người dùng tăng lên, PostgreSQL bắt đầu gặp giới hạn. Những server lớn sinh ra lượng message khổng lồ khiến một số bảng trở thành hot spot, truy vấn ngày càng chậm và việc mở rộng theo chiều dọc không còn hiệu quả.
Discord giải quyết bằng cách shard database theo Guild ID.
Những guild khác nhau sẽ được lưu trên các database khác nhau.
Điều này giúp workload được phân tán thay vì dồn vào một instance PostgreSQL duy nhất.
Redis chỉ lưu dữ liệu "nóng"
Không phải mọi dữ liệu đều cần đọc từ database. Một số dữ liệu thay đổi liên tục như:
Online status
Typing indicator
Presence
được Discord lưu trong Redis. Lý do rất đơn giản bởi đây là dữ liệu có vòng đời ngắn, được đọc rất nhiều nhưng không cần lưu trữ lâu dài.
Nếu mỗi lần người dùng bắt đầu gõ phím đều ghi xuống PostgreSQL thì database sẽ nhanh chóng trở thành bottleneck.
Redis giúp giảm tải rất lớn cho tầng lưu trữ chính và cải thiện đáng kể độ trễ.
Voice là phần khó nhất
So với chat, voice còn phức tạp hơn rất nhiều. Discord xây dựng hạ tầng VoIP riêng thay vì phụ thuộc hoàn toàn vào các giải pháp có sẵn.
Hệ thống sử dụng UDP để giảm latency, triển khai các cụm voice theo từng khu vực địa lý và bổ sung relay server để hỗ trợ người dùng nằm sau NAT hoặc firewall. Ngoài ra còn có heartbeat, tự động reconnect và cơ chế chọn region phù hợp nhằm đảm bảo chất lượng cuộc gọi.
Đây là lý do Discord thường có chất lượng voice rất ổn định ngay cả khi có hàng triệu người sử dụng đồng thời.
Quan sát hệ thống còn quan trọng hơn mở rộng hệ thống
Ở quy mô Internet, việc biết hệ thống đang hoạt động ra sao quan trọng không kém việc xây dựng nó. Discord sử dụng Prometheus và Grafana để thu thập metrics, kết hợp với Sentry để theo dõi lỗi ứng dụng. Họ còn xây dựng dashboard riêng để theo dõi sức khỏe của Gateway, tải của từng shard và thời gian phân phối sự kiện.
Khi một shard có dấu hiệu quá tải hoặc latency tăng bất thường, đội ngũ vận hành có thể phát hiện rất sớm trước khi người dùng nhận ra.
Kiến trúc cũng phải thay đổi theo thời gian
Một điểm rất đáng học là Discord không giữ nguyên stack ban đầu.
Khi workload ngày càng đa dạng, họ bắt đầu sử dụng thêm nhiều ngôn ngữ khác nhau.
Theo bài viết, những service yêu cầu hiệu năng cao dần được viết bằng Rust, Go và Kotlin. Hạ tầng cũng chuyển từ các bare-metal server được quản lý thủ công sang Kubernetes để tận dụng khả năng tự động mở rộng và orchestration.
Điều này phản ánh một tư duy rất thực tế:
Không có công nghệ nào phù hợp với mọi bài toán.
Kiến trúc nên thay đổi cùng với quy mô của sản phẩm.
Bài học lớn nhất từ Discord
Điều thú vị nhất trong câu chuyện của Discord không phải là họ sử dụng Elixir, PostgreSQL hay Redis. Bài học lớn hơn nằm ở cách họ tiếp cận việc mở rộng hệ thống. Họ không bắt đầu bằng Kubernetes, hàng chục microservice hay hàng trăm database.
Thay vào đó, Discord xây một hệ thống đủ đơn giản để ra mắt sản phẩm nhanh. Chỉ khi gặp bottleneck thực sự, họ mới bổ sung sharding, caching, voice cluster hay Kubernetes.
Đây cũng là lý do Discord có thể phát triển từ một ứng dụng chat dành cho game thủ thành một nền tảng phục vụ hàng triệu người dùng đồng thời mà không cần phải thay đổi toàn bộ hệ thống ngay từ đầu.
Kết luận
Kiến trúc realtime của Discord không dựa trên một công nghệ "thần kỳ". Thành công của họ đến từ việc kết hợp nhiều kỹ thuật đã được kiểm chứng:
Elixir/BEAM để xử lý hàng triệu kết nối đồng thời.
WebSocket Gateway để duy trì giao tiếp realtime.
Sharding để mở rộng theo chiều ngang và cô lập lỗi.
PostgreSQL Sharding để phân tán dữ liệu.
Redis cho dữ liệu thay đổi liên tục.
Voice clusters tối ưu cho độ trễ thấp.
Observability để theo dõi toàn bộ hệ thống theo thời gian thực.
Thông điệp lớn nhất của bài viết khá đơn giản nhưng rất đáng nhớ:
Đừng cố xây một hệ thống cho một tỷ người dùng ngay từ ngày đầu. Hãy xây một hệ thống đủ đơn giản để phát triển nhanh, sau đó mở rộng từng phần khi thực sự cần đến nó.