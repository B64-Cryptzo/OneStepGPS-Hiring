module.exports = {
    "env": {
        "browser": true,
        "es2021": true
    },
    "extends": [
        "eslint:recommended",
        "plugin:vue/essential"
    ],
    "parserOptions": {
        "ecmaVersion": 12,
        "sourceType": "module"
    },
    "plugins": [
        "vue"
    ],
    "rules": {
        "vue/valid-define-emits": "error",  // Enforce valid defineEmits usage
        "no-undef": "off"  // Turn off no-undef rule for Composition API
    }
};
