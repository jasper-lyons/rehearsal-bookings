
document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('booking-form');
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
        history.replaceState({ step: 1 }, null, ''); // Use replaceState to avoid adding a new history entry
        showStep(1);
    });

    // Back Button for Step 3
    document.getElementById('back-step-2').addEventListener('click', () => {
        history.replaceState({ step: 2 }, null, ''); // Use replaceState to avoid adding a new history entry
        showStep(2);
    });

    // Handle browser navigation
    window.addEventListener('popstate', (event) => {
        if (event.state && event.state.step) {
            showStep(event.state.step);
        }
    });

    // Initialize history state
    history.replaceState({ step: 1 }, null, '');

    // Validate Step 1
    const validateStep1 = () => {
        const sessionType = document.getElementById('session-type').value;
        const date = document.getElementById('date-input').value;
        const startTime = document.getElementById('start-time').value;
        return sessionType && date && startTime;
    };

    // Validate Step 2
    const validateStep2 = () => {
        const name = document.getElementById('name').value.trim();
        const email = document.getElementById('email').value.trim();
        const phone = document.getElementById('phone').value.trim();
       
        let nameRegex = /^[a-zA-Z]+([-' ][a-zA-Z]+)*\s+[a-zA-Z]+([-' ][a-zA-Z]+)*$/;
        let emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        let phoneRegex = /^\+?\d{11}$/;
      
        if (!nameRegex.test(name))  {
            alert('Please provide your full name.');
            return false;
        }
      
        if (!emailRegex.test(email)) {
            alert('Please enter a valid email address.');
            return false;
        }

        if (!phoneRegex.test(phone)) {
            alert('Please enter a valid phone number.');
            return false;
        }
        
        return name && email && phone;
    };
});		
