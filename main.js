document.addEventListener('DOMContentLoaded', () => {
    const navLinks = document.querySelectorAll('nav a');
    const sections = document.querySelectorAll('main section');
    const sidebar = document.querySelector('.sidebar');
    const examplesGrid = document.getElementById('examples-grid');
    const tagFilterButtons = document.querySelectorAll('.tag-filter-btn');

    // Function to show a specific section and update active link
    function showSection(id) {
        sections.forEach(section => {
            section.classList.remove('active');
        });
        const targetSection = document.getElementById(id);
        if (targetSection) {
            targetSection.classList.add('active');
            // Scroll to the section if not already in view
            targetSection.scrollIntoView({ behavior: 'smooth' });
        }

        navLinks.forEach(link => {
            link.classList.remove('active-link');
            if (link.getAttribute('data-section') === id) {
                link.classList.add('active-link');
            }
        });
    }

    // Handle navigation clicks
    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const sectionId = link.getAttribute('data-section');
            showSection(sectionId);
            // Update URL hash without triggering hashchange event immediately
            history.pushState(null, '', `#${sectionId}`);
        });
    });

    // Smooth scroll for internal links (e.g., from hero section)
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const targetId = this.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth'
                });
                history.pushState(null, '', `#${targetId}`);
                // Manually update active link for non-sidebar links
                navLinks.forEach(link => {
                    link.classList.remove('active-link');
                    if (link.getAttribute('data-section') === targetId) {
                        link.classList.add('active-link');
                    }
                });
            }
        });
    });


    // Intersection Observer for scroll spy
    const observerOptions = {
        root: null, // viewport
        rootMargin: '0px 0px -70% 0px', // Trigger when 30% of section is visible
        threshold: 0 // As soon as any part of the section enters the root
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const id = entry.target.id;
                navLinks.forEach(link => {
                    link.classList.remove('active-link');
                    if (link.getAttribute('data-section') === id) {
                        link.classList.add('active-link');
                    }
                });
            }
        });
    }, observerOptions);

    sections.forEach(section => {
        observer.observe(section);
    });

    // Initial load: show section based on URL hash or default to 'overview'
    const initialSectionId = window.location.hash ? window.location.hash.substring(1) : 'overview';
    showSection(initialSectionId);

    // Handle hash changes (e.g., back/forward browser buttons)
    window.addEventListener('hashchange', () => {
        const sectionId = window.location.hash ? window.location.hash.substring(1) : 'overview';
        showSection(sectionId);
    });

    // --- Examples Tag Filtering ---
    if (examplesGrid && tagFilterButtons.length > 0) {
        tagFilterButtons.forEach(button => {
            button.addEventListener('click', () => {
                const filterTag = button.getAttribute('data-tag');
                
                // Update active state of filter buttons
                tagFilterButtons.forEach(btn => btn.classList.remove('bg-indigo-600', 'text-white', 'hover:bg-indigo-700'));
                button.classList.add('bg-indigo-600', 'text-white', 'hover:bg-indigo-700');
                
                const exampleCards = examplesGrid.querySelectorAll('.example-card');
                exampleCards.forEach(card => {
                    const cardTags = card.getAttribute('data-tags').split(' ');
                    if (filterTag === 'all' || cardTags.includes(filterTag)) {
                        card.style.display = 'flex'; // Show card
                    } else {
                        card.style.display = 'none'; // Hide card
                    }
                });
            });
        });
    }
});
