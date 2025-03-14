
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
    };

    // Step 1 to Step 2
    toStep2Button.addEventListener('click', () => {
        if (validateStep1()) {
            populateSummary();
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
            showStep(3);
        }
    });

    // Back Button for Step 2
    document.getElementById('back-step-1').addEventListener('click', () => {
        showStep(1);
    });

    // Back Button for Step 3
    document.getElementById('back-step-2').addEventListener('click', () => {
        showStep(2);
    });

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

    const populateCustomerInfo = () => {
        // const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const phone = document.getElementById('phone').value;
        // const summaryName = document.querySelectorAll('.');
        const summaryEmail = document.getElementById('customer-email');
        const summaryPhone = document.getElementById('customer-phone');

        // summaryName.forEach(element => element.textContent = `Name: ${name}`);
        summaryEmail.innerHTML = `<p>An email confirmation will be sent to: <strong>${email}</strong></p>`;
        summaryPhone.innerHTML = `<p>On the day, an SMS with the access codes will be sent to: <strong>${phone}</strong></p>`;
    }
});		
