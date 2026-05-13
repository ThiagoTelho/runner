package br.gov.ses.hubsaude.assinador.model;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TipoMaterialTest {

    @Test
    void fromStringAcceptsExactName() {
        assertEquals(TipoMaterial.PEM, TipoMaterial.fromString("PEM"));
        assertEquals(TipoMaterial.PKCS12, TipoMaterial.fromString("PKCS12"));
        assertEquals(TipoMaterial.SMARTCARD, TipoMaterial.fromString("SMARTCARD"));
        assertEquals(TipoMaterial.TOKEN, TipoMaterial.fromString("TOKEN"));
        assertEquals(TipoMaterial.REMOTE, TipoMaterial.fromString("REMOTE"));
    }

    @Test
    void fromStringIsCaseInsensitive() {
        assertEquals(TipoMaterial.PEM, TipoMaterial.fromString("pem"));
        assertEquals(TipoMaterial.PKCS12, TipoMaterial.fromString("Pkcs12"));
        assertEquals(TipoMaterial.REMOTE, TipoMaterial.fromString("ReMoTe"));
    }

    @Test
    void fromStringReturnsNullForUnknown() {
        assertNull(TipoMaterial.fromString("BOGUS"));
        assertNull(TipoMaterial.fromString(""));
    }

    @Test
    void fromStringReturnsNullForNull() {
        assertNull(TipoMaterial.fromString(null));
    }

    @Test
    void valoresAceitosListsAllValues() {
        String s = TipoMaterial.valoresAceitos();
        for (TipoMaterial t : TipoMaterial.values()) {
            assertTrue(s.contains(t.name()), "Esperava encontrar " + t.name() + " em: " + s);
        }
    }
}
