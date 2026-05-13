package br.gov.ses.hubsaude.assinador.model;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class PoliticaRevogacaoTest {

    @Test
    void fromStringAcceptsCanonicalValues() {
        assertEquals(PoliticaRevogacao.WARN, PoliticaRevogacao.fromString("warn"));
        assertEquals(PoliticaRevogacao.SOFT_FAIL, PoliticaRevogacao.fromString("soft-fail"));
        assertEquals(PoliticaRevogacao.STRICT, PoliticaRevogacao.fromString("strict"));
    }

    @Test
    void fromStringIsCaseInsensitive() {
        assertEquals(PoliticaRevogacao.WARN, PoliticaRevogacao.fromString("WARN"));
        assertEquals(PoliticaRevogacao.SOFT_FAIL, PoliticaRevogacao.fromString("Soft-Fail"));
        assertEquals(PoliticaRevogacao.STRICT, PoliticaRevogacao.fromString("STRICT"));
    }

    @Test
    void fromStringRejectsEnumName() {
        assertNull(PoliticaRevogacao.fromString("SOFT_FAIL"));
    }

    @Test
    void fromStringReturnsNullForUnknown() {
        assertNull(PoliticaRevogacao.fromString("nao-existe"));
        assertNull(PoliticaRevogacao.fromString(""));
    }

    @Test
    void fromStringReturnsNullForNull() {
        assertNull(PoliticaRevogacao.fromString(null));
    }

    @Test
    void getValorReturnsCanonicalString() {
        assertEquals("warn", PoliticaRevogacao.WARN.getValor());
        assertEquals("soft-fail", PoliticaRevogacao.SOFT_FAIL.getValor());
        assertEquals("strict", PoliticaRevogacao.STRICT.getValor());
    }

    @Test
    void valoresAceitosListsCanonicalForm() {
        String s = PoliticaRevogacao.valoresAceitos();
        assertTrue(s.contains("warn"));
        assertTrue(s.contains("soft-fail"));
        assertTrue(s.contains("strict"));
    }
}
