INSERT INTO users (user_id, user_name, user_email, password_hash, role) 
VALUES 
    ('5f9d8b6f-40f1-4c27-b3e5-d4511ea6b48c', 'alice', 'alice@example.com', '$2a$10$xb9AluFpL.NSZJwrRk8jd.r/Ixx9CFRSaBQIpv0BungA28x/kgHAW', 'ADMIN'),--hashedpassword1
    ('1a23b4c5-6789-1011-a91b-22f56c0d7731', 'bob', 'bob@example.com', '$2a$10$IoltkpPiypjN3qWDfauJVedd0qgqKaMVDCApUk6kAmaGPPV0ascse', 'USER'),--hashedpassword2
    ('3e9e6899-4567-4e44-b82d-459cb5756d2c', 'carol', 'carol@example.com', '$2a$10$7bSCCnTF.sOtXSTkDYdImuwKzGEaNGLsYgD8hEoBV4767Rst/c9.a', 'USER'),--hashedpassword3
    ('7d3e690b-d29e-4704-a49f-3eae568877c0', 'dave', 'dave@example.com', '$2a$10$fnSULPmLmM2nLqLwlt1jZ.MeJR2lU8aWClTjG1GEQZWmP86ohAxpS', 'USER'),--hashedpassword4
    ('9d44e54b-8701-423d-a06f-9b7bba63e768', 'eve', 'eve@example.com', '$2a$10$FJKMeg/FIyGqnX66IxzBzud3MwCx1Zm83as.sTLJQP6/rWoRFxRjO', 'ADMIN')--hashedpassword5
ON CONFLICT DO NOTHING;