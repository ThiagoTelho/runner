package br.gov.ses.hubsaude.assinador.model;

import java.util.Arrays;
import java.util.stream.Collectors;

public enum PoliticaRevogacao {
    WARN("warn"),
    SOFT_FAIL("soft-fail"),
    STRICT("strict");

    private final String valor;

    PoliticaRevogacao(String valor) {
        this.valor = valor;
    }

    public String getValor() {
        return valor;
    }

    public static PoliticaRevogacao fromString(String value) {
        if (value == null) return null;
        for (PoliticaRevogacao p : values()) {
            if (p.valor.equalsIgnoreCase(value)) return p;
        }
        return null;
    }

    public static String valoresAceitos() {
        return Arrays.stream(values())
                .map(PoliticaRevogacao::getValor)
                .collect(Collectors.joining(", "));
    }
}
