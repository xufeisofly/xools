;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;                            ;;
;;       multipath-merge      ;;
;;                            ;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; 目前只实现 L > K，即每一路 list 的 length > topK 的 K 值

(define l1 '(9 7 5 4 2 1))
(define l2 '(12 8 7 5 3 1))
(define l3 '(11 4 3 2 1))
(define l4 '(23 8 6 4 3 2 1))
(define l5 '(30 20 4 3 2 1))
(define l6 '(33 32 22 16 4 9 6))
(define l7 '(39 22 21 14 3 2))
(define l8 '(27 21 18 17 6 5))

(define list-bundle (list l1 l2 l3 l4 l5 l6 l7 l8))

;; 初始化一个各数组指针位置的列表，初始值 '(0 0 0 ...)
(define (init-idx-l lists)
  (define (iter result n)
    (if (= (length lists) n)
        result
        (iter (append result '(0)) (+ n 1))))
  (iter '() 0))


(define lists-ptrs (init-idx-l list-bundle))

;; node constructor 因为多路归并是二维数组，所以构建 item，包含数组的 list-n 和 value
(define (make-item l-n value)
  (cons l-n value))

(define (make-node index item)
  (cons index item))

(define (get-index node)
  (car node))

(define (get-value node)
  (cdr (cdr node)))

(define (get-item node)
  (cdr node))

;; node 的左子叶节点
(define (left-child-node l node)
  (let ((index (+ (* 2 (get-index node)) 1)))
    (node-by-index l index)))

;; node 的右子叶节点
(define (right-child-node l node)
  (if (and (even? (length l)) (= (get-index node) (get-index (last-parent-node l))))
      '()
      (let ((index (+ (* 2 (get-index node)) 2)))
        (node-by-index l index))))

(define (node-by-index l index)
  (make-node index (list-ref l index)))

(define (even? num)
  (= (remainder num 2) 0))

;; list 最后一个父节点
(define (last-parent-node l)
  (let ((len (length l)))
    (if (even? len)
        (make-node (- (/ len 2) 1) (list-ref l (- (/ len 2) 1)))
        (make-node (- (/ (- len 1) 2) 1) (list-ref l (- (/ (+ len 1) 2) 1))))))


;; 设置 list k v
(define (list-set! list k val)
  (if (zero? k)
      (set-car! list val)
      (list-set! (cdr list) (- k 1) val)))

;; 比较两个节点
(define (smaller? node1 node2)
  (< (get-value node1) (get-value node2)))

;; 取更大的节点
(define (bigger-node node1 node2)
  (cond ((null? node1) node2)
        ((null? node2) node1)
        (else
         (if (smaller? node1 node2)
             node2
             node1))))

;; 返回大的子节点
(define (bigger-child-node l p-node)
  (let ((node1 (left-child-node l p-node)) (node2 (right-child-node l p-node)))
    (bigger-node node1 node2)))

;; 是否根节点
(define (root-node? node)
  (= (get-index node) 0))

;; node 的父节点
(define (parent-node l node)
  (if (root-node? node)
      node
      (let ((node-index (get-index node)))
        (if (even? node-index)
            (node-by-index l (- (/ node-index 2) 1))
            (node-by-index l (- (/ (+ node-index 1) 2) 1))))))

;; 交换 l 的两个 node, l 被修改
(define (swap l node p-node)
  (begin
    (set! tmp (get-item node))
    (list-set! l (get-index node) (get-item p-node))
    (list-set! l (get-index p-node) tmp)))

(define (parent-node? l node)
  (<= (get-index node) (get-index (last-parent-node l))))

;; 子叶节点递归上浮
(define (flow-up l node)
  (let ((p-node (parent-node l node)))
    (if (smaller? p-node (bigger-child-node l p-node))
        (begin
          (let ((c-node (bigger-child-node l p-node)))
            (swap l p-node c-node)
            ;; 上浮后子树做下沉，微调
            (flow-down l (node-by-index l (get-index c-node)))
            (flow-up l (node-by-index l (get-index p-node))))))))

(define (flow-down l node)
  (if (parent-node? l node)
      (let ((c-node (bigger-child-node l node)))
        (if (smaller? node c-node)
            (begin
              (swap l node c-node)
              (flow-down l (node-by-index l (get-index c-node))))))))

(define (swap-with-child? l p-node)
  (let ((bigger-child (bigger-child-node l p-node)))
    (smaller? p-node bigger-child)))

;; 建堆，最大堆
(define (make-max-heap l)
  (define (iter p-idx)
    (if (< p-idx 0)
        l
        (let ((p-node (node-by-index l p-idx)))
          (if (swap-with-child? l p-node)
              (let ((up-node (bigger-child-node l p-node)))
                (flow-up l up-node)
                (iter (- p-idx 1)))
            (iter (- p-idx 1))))))
  (iter (get-index (last-parent-node l))))

(define get-n-items
  (lambda (lst num)
    (if (> num 0)
        (cons (car lst) (get-n-items (cdr lst) (- num 1)))
        '()))) ;'

(define slice
  (lambda (lst start count)
    (if (> start 1)
        (slice (cdr lst) (- start 1) count)
        (get-n-items lst count))))

(define (root-node h)
  (node-by-index h 0))

(define (last-node h)
  (node-by-index h (- (length h) 1)))

(define (delete-root-node h)
  (swap h (root-node h) (last-node h))
  (remove! (get-item (last-node h)) h)
  (flow-down h (root-node h)))

;; 返回堆的最大值，并 push 新值进堆， 会改变 h 值
(define (pop-max-item h lists)
  (set! select-item (get-item (root-node h)))
  (delete-root-node h)
  (push-item (next-item select-item lists) h)
  select-item)

;; 插入 item 到最大堆
(define (push-item item h)
  (append! h (list item))
  (make-max-heap h))

;; 取 item 的 list-n
(define (get-l-n item)
  (car item))


;; 获取 list-n 的下一个 item
(define (next-item select-item lists)
  (let ((l-n (get-l-n select-item)))
    (let ((cur-l (list-ref lists l-n)))
      (let ((cur-idx (list-ref lists-ptrs l-n)))
        (if (= cur-idx (- (length cur-l) 1))
            (error "list length should > K")
            (begin
              (list-set! lists-ptrs l-n (+ 1 cur-idx))
              (make-item l-n (list-ref cur-l (+ 1 cur-idx)))))))))


;; 用多路 list 的首个元素组成一个 list
(define (get-first-items lists)
  (define (iter result idx)
    (if (= (length lists) (+ idx))
        result
        (iter (append result
                      (list (make-item idx (car (list-ref lists idx)))))
              (+ idx 1))))
  (iter '() 0))


;; 多路归并！！！
(define (topK k lists)

  ;; 用首个元素 list 建堆
  (define (init-max-heap lists)
    (make-max-heap (get-first-items lists)))

  (define (topK-items k lists)
    (define max-h (init-max-heap lists))
    (define (iter result n)
      (if (= n k)
          result
          (iter (append result
                        (list (pop-max-item max-h lists)))
                (+ 1 n))))
    (iter '() 0))

  (define tK-items (topK-items k lists))

  (define (iter result n)
    (if (= n k)
        result
        (iter (append result
                      (list (cdr (list-ref tK-items n))))
              (+ n 1))))
  (iter '() 0))

;; (topK 10 list-bundle)

