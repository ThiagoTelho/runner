package br.gov.ses.hubsaude.assinador.model;

import java.util.Arrays;
import java.util.stream.Collectors;

public enum TipoMaterial {
    PEM,
    PKCS12,
    SMARTCARD,
    TOKEN,
    REMOTE;

    public static TipoMaterial fromString(String value) {
        if (value == null) return null;
        try {
            return valueOf(value.toUpperCase());
        } catch (IllegalArgumentException e) {
            return null;
        }
    }

    public static String valoresAceitos() {
        return Arrays.stream(values())
                .map(Enum::name)
                .collect(Collectors.joining(", "));
    }
}
