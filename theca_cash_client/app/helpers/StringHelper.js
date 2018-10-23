export class StringHelper {
    static isNullOrEmpty = (str) => {
        if (str === null || str === undefined) {
            return true;
        }
        return StringHelper.trim(str).length === 0;
    };

    static trim = (str, chars) => {
        return StringHelper.leftTrim(StringHelper.rightTrim(str, chars), chars);
    };

    static leftTrim = (str, charsParam) => {
        const chars = charsParam || '\\s';
        return str.replace(new RegExp('^[' + chars + ']+', 'g'), '');
    };

    static rightTrim = (str, charsParam) => {
        const chars = charsParam || '\\s';
        return str.replace(new RegExp('[' + chars + ']+$', 'g'), '');
    };
}
