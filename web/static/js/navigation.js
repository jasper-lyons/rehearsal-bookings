
document.addEventListener('DOMContentLoaded', () => {
    const steps = {
        1: document.getElementById('step-1'),
        2: document.getElementById('step-2'),
        3: document.getElementById('step-3'),
    };

    // Buttons for navigation
    const toStep2Button = document.getElementById('book-now');
    const toStep3Button = document.getElementById('pay-now');

    // Show/Hide steps
    const showStep = (stepNumber) => {
        Object.values(steps).forEach((step) => step.classList.add('hidden'));
        steps[stepNumber].classList.remove('hidden');
        setTimeout(() => {
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }, 0); // Use a timeout to ensure it overrides the browser's default behavior
    };

    // Step 1 to Step 2
    toStep2Button.addEventListener('click', () => {
        if (validateStep1()) {
            populateSummary();
            history.pushState({ step: 2 }, null, '');
            showStep(2);
        } else {
            alert('Please select a time!');
        }
    });

    // Step 2 to Step 3
    toStep3Button.addEventListener('click', () => {
        if (validateStep2()) {
            populateSummary();
            populateCustomerInfo();
            history.pushState({ step: 3 }, null, '');
            showStep(3);
        }
    });

    // Back Button for Step 2
    document.getElementById('back-step-1').addEventListener('click', () => {
        history.back();
        showStep(1);
    });

    // Back Button for Step 3
    document.getElementById('back-step-2').addEventListener('click', () => {
        history.back();
        showStep(2);
    });

    // Handle browser navigation
    window.addEventListener('popstate', (event) => {
        if (event.state && event.state.step) {
            showStep(event.state.step);
        }
    });

    // Initialize history state
    history.pushState({ step: 1 }, null, '');

    window.onbeforeunload = function() {
        return "Your rehearsal details will be lost if you leave the page, are you sure?";
    };

    // Validate Step 1
    const validateStep1 = () => {
        const sessionType = document.getElementById('session-type').value;
        const date = document.getElementById('date-input').value;
        const startTime = document.getElementById('start-time').value;
        return sessionType && date && startTime;
    };

    // Validate Step 2
    const validateStep2 = () => {
        const name = document.getElementById('name');
        const email = document.getElementById('email');
        const phone = document.getElementById('phone');

        let nameRegex = /^[a-zA-Z]+([-' ][a-zA-Z]+)*\s+[a-zA-Z]+([-' ][a-zA-Z]+)*$/;
        let emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        let phoneRegex = /^(?:\+44|44|0)?\s?7\d{3}\s?\d{3}\s?\d{3}$/;

        let isValid = true;

        // Reset previous error states
        [name, email, phone].forEach((input) => {
            input.classList.remove('error');
            const errorSpan = input.nextElementSibling;
            if (errorSpan && errorSpan.classList.contains('error-message')) {
                errorSpan.remove();
            }
        });

        if (!nameRegex.test(name.value.trim())) {
            isValid = false;
            name.classList.add('error');
            const errorSpan = document.createElement('span');
            errorSpan.classList.add('error-message');
            errorSpan.textContent = 'Please provide your full name.';
            name.parentNode.appendChild(errorSpan);
        }

        if (!emailRegex.test(email.value.trim())) {
            isValid = false;
            email.classList.add('error');
            const errorSpan = document.createElement('span');
            errorSpan.classList.add('error-message');
            errorSpan.textContent = 'Please enter a valid email address.';
            email.parentNode.appendChild(errorSpan);
        }

        if (!phoneRegex.test(phone.value.trim())) {
            isValid = false;
            phone.classList.add('error');
            const errorSpan = document.createElement('span');
            errorSpan.classList.add('error-message');
            errorSpan.textContent = 'Please enter a valid phone number.';
            phone.parentNode.appendChild(errorSpan);
        }

        if (!isValid) {
            alert('Please enter all your details.');
        }

        return isValid;
    };
});		
