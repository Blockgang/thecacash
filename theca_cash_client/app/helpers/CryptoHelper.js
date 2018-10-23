const bcrypt = require('bcryptjs');

export function createPasswordHash(password) {
    const salt = bcrypt.genSaltSync(13);
    return bcrypt.hashSync(password, salt);
}

export function comparePasswordWithHash(password, hash) {
    return bcrypt.compareSync(password, hash);
}

