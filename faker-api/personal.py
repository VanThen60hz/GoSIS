import requests
import json
from faker import Faker

fake = Faker('vi_VN')  # Vietnamese locale

url = "http://localhost:6000/personal"
headers = {'Content-Type': 'application/json'}

n = int(input("Nhập số lượng request : "))

for i in range(n):
    payload = json.dumps({
        "Employee_ID": fake.random_number(digits=4, fix_len=True),
        "First_Name": fake.first_name(),
        "Last_Name": fake.last_name(),
        "Middle_Initial": fake.random_letter(),
        "Address1": fake.street_address(),
        "Address2": fake.secondary_address(),
        "City": fake.city(),
        "State": fake.state(),
        "Zip": int(fake.zipcode()),
        "Email": fake.email(),
        "PhoneNumber": fake.phone_number(),
        "SocialSecurityNumber": fake.ssn(),
        "DriversLicense": fake.license_plate(),
        "Marital_Status": fake.random_element(elements=("Single", "Married", "Divorced")),
        "Gender": fake.random_element(elements=(True, False)),
        "Shareholder_Status": fake.random_element(elements=(True, False)),
        "Benefit_Plans": fake.random_int(min=1, max=3),
        "Ethnicity": "Kinh"
    })

    response = requests.request("POST", url, headers=headers, data=payload)

    if response.status_code == 200:
        print(f'POST request {i+1} successful')
    else:
        print(f'Error on POST request {i+1}: {response.text}')
