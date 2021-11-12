const warning = [
  {
    selector: "[accesskey]",
    category: "warning",
    name: "accesskey-present",
    description:
      "Accesskey attribute could interfere and conflict with screen readers and assistive technologies",
  },
  {
    selector:
      "[aria-hidden=true] a, [aria-hidden=true] button, [aria-hidden=true] summary, [aria-hidden=true] area, [aria-hidden=true] audio[controls], [aria-hidden=true] iframe, [aria-hidden=true] input, [aria-hidden=true] select, [aria-hidden=true] textarea, [aria-hidden=true] video[controls]",
    category: "warning",
    name: "focusable-element-in-aria-hidden-element",
    description: "Focusable element in a aria-hidden element",
  },
  {
    selector: "[autoplay]",
    category: "warning",
    name: "autoplay-present",
    description: "Media should not autoplay",
  },
  {
    selector:
      "[onclick]:not(a):not(button):not([onkeydown]):not([onkeyup]):not([onkeypress]), [onclick][onkeyup]:not([tabindex]):not(a):not(button), [onclick][onkeyup][tabindex='-1']:not(a):not(button), [onclick][onkeydown]:not([tabindex]):not(a):not(button), [onclick][onkeydown][tabindex='-1']:not(a):not(button), [onclick][onkeypress]:not([tabindex]):not(a):not(button), [onclick][onkeypress][tabindex='-1']:not(a):not(button), [ondblclick]:not(a):not(button):not([onkeydown]):not([onkeyup]):not([onkeypress]), [ondblclick][onkeyup]:not([tabindex]):not(a):not(button), [ondblclick][onkeyup][tabindex='-1']:not(a):not(button), [ondblclick][onkeydown]:not([tabindex]):not(a):not(button), [ondblclick][onkeydown][tabindex='-1']:not(a):not(button), [ondblclick][onkeypress]:not([tabindex]):not(a):not(button), [ondblclick][onkeypress][tabindex='-1']:not(a):not(button), [onmousedown]:not(a):not(button):not([onkeydown]):not([onkeyup]):not([onkeypress]), [onmousedown][onkeyup]:not([tabindex]):not(a):not(button), [onmousedown][onkeyup][tabindex='-1']:not(a):not(button), [onmousedown][onkeydown]:not([tabindex]):not(a):not(button), [onmousedown][onkeydown][tabindex='-1']:not(a):not(button), [onmousedown][onkeypress]:not([tabindex]):not(a):not(button), [onmousedown][onkeypress][tabindex='-1']:not(a):not(button), [onmouseup]:not(a):not(button):not([onkeydown]):not([onkeyup]):not([onkeypress]), [onmouseup][onkeyup]:not([tabindex]):not(a):not(button), [onmouseup][onkeyup][tabindex='-1']:not(a):not(button), [onmouseup][onkeydown]:not([tabindex]):not(a):not(button), [onmouseup][onkeydown][tabindex='-1']:not(a):not(button), [onmouseup][onkeypress]:not([tabindex]):not(a):not(button), [onmouseup][onkeypress][tabindex='-1']:not(a):not(button), [onmouseover]:not(a):not(button), [onmouseenter]:not(a):not(button), [onmouseleave]:not(a):not(button)",
    category: "warning",
    name: "inaccessible-click-event-on-non-clickable-element",
    description:
      "Potentially inaccessible click event used on non-clickable element. Ensure that you have an accessible alternative",
  },
  {
    selector:
      "h1[role=text], h2[role=text], h3[role=text], h4[role=text], h5[role=text], h6[role=text]",
    category: "warning",
    name: "text-role-on-heading-element",
    description:
      "Using role='text' on a heading element causes it to lose semantic meaning for screen readers",
  },
  {
    selector:
      "img[alt^='logo for' i], img[alt^='logo of' i], img[alt^='image of' i], img[alt^=bullet i], img[alt^='graphic of' i], img[alt$=logo i], img[alt$=graphic i], img[alt$=image i], img[alt$='.ivaf' i], img[alt$='.webp' i], img[alt$='.bmp' i], img[alt$='.svg' i], img[alt$=',jpeg' i], img[alt$='.jpg' i], img[alt$='.gif' i], img[alt$='.png' i], img[alt=' ' i], img[alt='*' i], img[alt=graph i], img[alt=diagram i], img[alt=table i], img[alt=chart i], img[alt=blank i], img[alt=spacer i], img[alt=more i], img[alt=arrow i], img[alt=button i], img[alt=bullet i], img[alt=logo i], img[alt=artwork i], img[alt=painting i], img[alt=drawing i], img[alt=photograph i], img[alt=photo i], img[alt=graphic i], img[alt=image i]",
    category: "warning",
    name: "nondescriptive-alt-text",
    description: "Alt text must describe the content of the image",
  },
  {
    selector: "[style*='!important']",
    category: "warning",
    name: "nondescriptive-alt-text",
    description: "!important inline styles should be avoided",
  },
  {
    selector: "*:not([role=text]) > em, *:not([role=text]) span",
    category: "warning",
    name: "missing-text-role-in-broken-text",
    description:
      "Text broken by <span> or <em> elements might need a [role='text'] attribute",
  },
  {
    selector: "a[target=_blank]",
    category: "warning",
    name: "new-tab-links",
    description:
      "Links that open in a new tab or window should be communicated to the user, ensure the user knows of this behavior",
  },
  {
    selector: "a[href='#'], a[role=button], a[href^='javascript:']",
    category: "warning",
    name: "anchor-tag-as-button",
    description:
      "Anchor tags should not be used as buttons. Links should redirect to a resource/page, if they don't they probably should be buttons",
  },
  {
    selector:
      "a[href$='.pdf'], a[href$='.docx'], a[href$='.doc'], a[href$='.xlsx'], a[href$='.xls'], a[href$='.pptx'], a[href$='.pptm'], a[href$='.ppt'], a[href$='.txt']",
    category: "warning",
    name: "link-to-document",
    description:
      "Anchor tags that contain a link to a PDF, Word, Excel, or PowerPoint document is present on the page and could be a potential accessibility issue",
  },
  {
    selector:
      "a[href*='https://youtu.be'], a[href*='https://youtube.com'], a[href*='https://spotify.com/']",
    category: "warning",
    name: "link-to-content-site",
    description:
      "Anchor tags that contain a link to content site (Youtube, Spotify etc) is present on the document.",
  },
  {
    selector: "section:empty, section > :not(h1, h2, h3, h4, h5, h6, img)",
    category: "warning",
    name: "missing-section-heading",
    description: "A <section> should contain a heading element",
  },
  {
    selector: "[title]",
    category: "warning",
    name: "title-present",
    description:
      "The title attribute has many accessibility concerns, consider using an alternative method",
  },
  {
    selector: "u",
    category: "warning",
    name: "u-present",
    description:
      "Using <u> element could be confused as a hyperlink, consider using a different element such as <em> or <b>",
  },
];

module.exports = warning;
