;; emacs mode for OldRope
;; does no work for some reason
(setq oldrope-directive-keywords-regexp (regexp-opt  '("page" "link" "goto" "act" "end" "div" "span" "include")))
 (defvar oldrope-font-lock-defaults
         `(((,oldrope-directive-keywords-regexp . font-lock-keyword-face))))
(define-derived-mode oldrope-mode fundamental-mode
  "oldrope mode"
  "Major mode for editing OldRope games"
  (setq comment-start "/*")
  (setq comment-end "*/")
  (setq font-lock-defaults oldrope-font-lock-defaults))
(provide 'oldrope-mode)
