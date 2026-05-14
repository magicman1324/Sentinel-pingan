package com.pingan.monitor.deserializer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.pingan.monitor.model.MetricEvent;
import org.apache.flink.api.common.serialization.DeserializationSchema;
import org.apache.flink.api.common.typeinfo.TypeInformation;

import java.io.IOException;

public class MetricDeserializer implements DeserializationSchema<MetricEvent> {

    private static final ObjectMapper mapper = new ObjectMapper();

    @Override
    public MetricEvent deserialize(byte[] message) throws IOException {
        return mapper.readValue(message, MetricEvent.class);
    }

    @Override
    public boolean isEndOfStream(MetricEvent nextElement) {
        return false;
    }

    @Override
    public TypeInformation<MetricEvent> getProducedType() {
        return TypeInformation.of(MetricEvent.class);
    }
}
