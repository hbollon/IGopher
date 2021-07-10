import * as config from "@/config";

config.ready(() => {
  // Toggle the side navigation
  const sidebar = document.querySelector(".sidebar");
  const sidebarToggles = document.querySelectorAll(
    "#sidebarToggle, #sidebarToggleTop"
  );

  if (sidebar) {
    const collapseEl = sidebar.querySelector(".collapse");
    const collapseElementList = [].slice.call(
      document.querySelectorAll(".sidebar .collapse")
    );
    const sidebarCollapseList = collapseElementList.map(function(collapseEl) {
      return new config.bootstrap.Collapse(collapseEl, { toggle: false });
    });

    for (const toggle of sidebarToggles) {
      // Toggle the side navigation
      toggle.addEventListener("click", function(e) {
        document.body.classList.toggle("sidebar-toggled");
        sidebar.classList.toggle("toggled");

        if (sidebar.classList.contains("toggled")) {
          for (const bsCollapse of sidebarCollapseList) {
            bsCollapse.hide();
          }
        }
      });
    }

    // Close any open menu accordions when window is resized below 768px
    window.addEventListener("resize", function() {
      const vw = Math.max(
        document.documentElement.clientWidth || 0,
        window.innerWidth || 0
      );
      if (vw < 768) {
        for (const bsCollapse of sidebarCollapseList) {
          bsCollapse.hide();
        }
      }
    });
  }

  // Prevent the content wrapper from scrolling when the fixed side navigation hovered over
  const fixedNav = document.querySelector("body.fixed-nav .sidebar");
  if (fixedNav) {
    fixedNav.addEventListener("mousewheel DOMMouseScroll wheel", function(e) {
      const vw = Math.max(
        document.documentElement.clientWidth || 0,
        window.innerWidth || 0
      );
      if (vw > 768) {
        // let delta = e.wheelDelta || -e.detail;
        // document.body.scrollTop += (delta < 0 ? 1 : -1) * 30;
        // e.preventDefault();
      }
    });
  }

  // Scroll to top button appear
  const scrollToTop = document.querySelector(".back-to-top") as HTMLElement;
  window.addEventListener("scroll", function() {
    const scrollDistance = window.pageYOffset;

    // Check if user is scrolling up
    if (scrollToTop != null) {
      if (scrollDistance > 100) {
        scrollToTop.style.display = "block";
      } else {
        scrollToTop.style.display = "none";
      }
    }
  });

  // Scroll to top button callback
  function backToTop() {
    const rootElement = document.documentElement
    rootElement.scrollTo({
      top: 0,
      behavior: "smooth"
    })
  }
  scrollToTop.addEventListener("click", backToTop);

});
