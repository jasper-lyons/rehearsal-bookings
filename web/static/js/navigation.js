
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
            showStep(3);
            
        } else {
            alert('Please fill in all details.');
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
        return name && email && phone;
    };

    // Populate Summary for Step 3
    const populateSummary = () => {
        const room = document.getElementById('room').value;
        const date = document.getElementById('date-input').value;
        const startTime = document.getElementById('start-time').value;
        const endTime = document.getElementById('end-time').value;
        const price = document.getElementById('price').textContent;

        // Update both Step 2 and Step 3 summaries using classes
        const summaryRooms = document.querySelectorAll('.summary-room');
        const summaryDate= document.querySelectorAll('.summary-date');
        const summaryTimes = document.querySelectorAll('.summary-time');
        const summaryPrices = document.querySelectorAll('.summary-price');

        // Update the summary content for all elements with these classes
        summaryRooms.forEach(element => element.textContent = `Room: ${room}`);
        summaryDate.forEach(element => element.textContent = `Date: ${date}`);
        summaryTimes.forEach(element => element.textContent = `Time: ${startTime} - ${endTime}`);
        summaryPrices.forEach(element => element.textContent = `Price: ${price}`);
    };
});		
